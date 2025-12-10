package taskvalidator

import (
	"context"
	"time"
	"todoapp/internal/param"
	"todoapp/internal/pkg/errmsg"
	"todoapp/internal/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidateCreateTask(ctx context.Context ,req param.CreateTaskRequest)(map[string]string, error){
	const op=richerror.Op("taskvalidator.ValidateCreateTask")

	err:=validation.ValidateStruct(&req,
		validation.Field(&req.Title,
		validation.Required.Error("title is required"),
		validation.Length(3,50).Error("title must be 3-50 chracter long")),

		validation.Field(&req.CategoryID,
			validation.Required.Error("categoryis required"),
        ),

		 validation.Field(&req.Priority,
            validation.Required.Error("priority is required"),
            validation.In("low", "medium", "high").Error("priority must be low, medium, or high")),
		    validation.Field(&req.Description,
            validation.By(v.validateOptionalDescription),
        ),
	)

    fieldErrors:=make(map[string]string)

	if err !=nil {
        if errV,ok:=err.(validation.Errors);ok{
            for key,value :=range errV {
                if value!=nil{
                    fieldErrors[key]=value.Error()
                }
            }
        }
    }

    if req.CategoryID >0{
        _,err:=v.categoryRepo.GetCategoryByID(ctx,req.CategoryID,req.UserID)
        if err != nil {
			fieldErrors["category_id"] = "category not found or you don't have permission"
		}
    }

    if req.DueDate != nil && req.DueDate.Before(time.Now()) {
		fieldErrors["due_date"] = "due date cannot be in the past"
	}

    	if len(fieldErrors) > 0 {
		return fieldErrors, richerror.New(op).
			WithMessage(errmsg.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta("fields", fieldErrors)
	}

	return nil, nil

	
}







