package taskvalidator

import (
	"context"
	"fmt"
	"strings"
	"time"
	"todoapp/internal/param"
	"todoapp/internal/pkg/errmsg"
	"todoapp/internal/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidateUpdateTask(ctx context.Context, taskID uint, req param.UpdateTaskRequest) (map[string]string, error) {
	const op = richerror.Op("taskvalidator.ValidateUpdateTask")

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Title,
			validation.By(v.validateOptionalTitle),
		),
		validation.Field(&req.Description,
			validation.By(v.validateOptionalDescription),
		),
		validation.Field(&req.Priority,
			validation.By(v.validateOptionalPriority),
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

	// Check category if being updated
	if req.CategoryID != nil && *req.CategoryID > 0 {
		_, err := v.categoryRepo.GetCategoryByID(ctx, *req.CategoryID, req.UserID)
		if err != nil {
			fieldErrors["category_id"] = "category not found or you don't have permission"
		}
	}

	// Validate due date if provided
	if req.DueDate != nil && req.DueDate.Before(time.Now()) {
		fieldErrors["due_date"] = "due date cannot be in the past"
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

func (v Validator) validateOptionalTitle(value any)error {
	if value ==nil{
		return  nil
	}

	title,ok:=value.(*string)
	if !ok ||title ==nil{
 		return nil
	}

	titleLen:=len(*title)
	if titleLen<3 || titleLen>50 {
		return fmt.Errorf("title must be 3-50 chracter long")
	}
	return nil
}

func (v Validator) validateOptionalPriority(value interface{}) error {
	if value == nil {
		return nil
	}

	priority, ok := value.(*string)
	if !ok || priority == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*priority)

	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}

	if trimmed != "" && !validPriorities[trimmed] {
		return fmt.Errorf("priority must be low, medium, or high")
	}

	return nil
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
    
    // If they sent it, make sure it's meaningful
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