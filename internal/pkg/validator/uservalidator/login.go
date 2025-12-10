package uservalidator

import (
	"context"
	"errors"
	"todoapp/internal/param"
	"todoapp/internal/pkg/errmsg"
	"todoapp/internal/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidateLoginRequest(ctx context.Context, req param.LoginRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateLoginRequest"

	err := validation.ValidateStruct(&req,
		validation.Field(&req.PhoneNumber,
			validation.Required,
			validation.Match(phoneNumberRegex).Error("phone number is not valid"),
		),
		validation.Field(&req.Password,
			validation.Required.Error("password is required"),
			validation.Length(8, 0).Error("password must be at least 8 characters long"),
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

	// Check if phone number exists with context
	if req.PhoneNumber != "" {
		_, err := v.repo.GetUserByPhoneNumber(ctx, req.PhoneNumber)
		if err != nil {
			// Check if it's a "not found" error
			var re *richerror.RichError
			if errors.As(err, &re) && re.GetKind() == richerror.KindNotFound {
				fieldErrors["phone_number"] = "phone number not found"
			} else {
				// Unexpected database error
				return fieldErrors, richerror.New(op).
					WithMessage("failed to check phone number").
					WithKind(richerror.KindUnexpected).
					WithErr(err)
			}
		}
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

