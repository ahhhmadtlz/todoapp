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

func (v Validator) ValidateUpdateCategory(ctx context.Context, categoryID uint, req param.UpdateCategoryRequest) (map[string]string, error) {
	const op = richerror.Op("categoryvalidator.ValidateUpdateCategory")

	// Basic validation
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.By(v.validateOptionalName),
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

	// Check uniqueness if name is being updated
	if req.Name != nil && *req.Name != "" {
		trimmed := strings.TrimSpace(*req.Name)
		existingCat, err := v.repo.GetCategoryByName(ctx, req.UserID, trimmed)
		
		if err == nil {
			// Found a category with this name
			// Check if it's a DIFFERENT category
			if existingCat.ID != categoryID {
				fieldErrors["name"] = "category name already exists"
			}
			// If same ID, it's OK (user didn't change name or changed to same value)
		}
		// If error (not found), name is unique - OK
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

func (v Validator) validateOptionalName(value interface{}) error {
	if value == nil {
		return nil
	}

	name, ok := value.(*string)
	if !ok || name == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*name)

	if trimmed == "" {
		return fmt.Errorf("category name cannot be empty or whitespace only")
	}

	nameLen := len(trimmed)
	if nameLen < 3 {
		return fmt.Errorf("category name must be at least 3 characters long")
	}

	if nameLen > 50 {
		return fmt.Errorf("category name must be less than 50 characters")
	}

	return nil
}