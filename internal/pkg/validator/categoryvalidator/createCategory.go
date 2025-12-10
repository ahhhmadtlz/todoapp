package categoryvalidator

import (
	"context"
	"fmt"
	"strings"
	"todoapp/internal/param"
	"todoapp/internal/pkg/errmsg"
	"todoapp/internal/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidateCreateCategory(ctx context.Context, req param.CreateCategoryRequest) (map[string]string, error) {
	const op = richerror.Op("categoryvalidator.ValidateCreateCategory")

	// Basic validation first (no DB calls)
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.Required.Error("category name is required"),
			validation.Length(3, 50).Error("category name must be 3-50 characters long"),
		),
		validation.Field(&req.Description,
			validation.By(v.validateOptionalDescription),
		),
	)

	fieldErrors := make(map[string]string)

	// Collect validation errors
	if err != nil {
		if errV, ok := err.(validation.Errors); ok {
			for key, value := range errV {
				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}
	}

	// Check uniqueness separately (needs context)
	if req.Name != "" {
		trimmed := strings.TrimSpace(req.Name)
		_, err := v.repo.GetCategoryByName(ctx, req.UserID, trimmed)
		if err == nil {
			// Category found - name is NOT unique
			fieldErrors["name"] = "category name already exists"
		}
		// If error (not found), name is unique - good!
	}

	// Return errors if any
	if len(fieldErrors) > 0 {
		return fieldErrors, richerror.New(op).
			WithMessage(errmsg.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta("fields", fieldErrors)
	}

	return nil, nil
}

func (v Validator) validateOptionalDescription(value interface{}) error {
	if value == nil {
		return nil
	}

	desc, ok := value.(*string)
	if !ok || desc == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*desc)

	if trimmed == "" {
		return fmt.Errorf("description cannot be empty or whitespace only")
	}

	if len(trimmed) < 3 {
		return fmt.Errorf("description must be at least 3 characters long")
	}

	if len(trimmed) > 500 {
		return fmt.Errorf("description must be less than 500 characters")
	}

	return nil
}