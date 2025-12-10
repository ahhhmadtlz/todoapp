package uservalidator

import (
	"context"
	"todoapp/internal/param"
	"todoapp/internal/pkg/errmsg"
	"todoapp/internal/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
)



func (v Validator) ValidateRegisterRequest(ctx context.Context, req param.RegisterRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateRegisterRequest"

	// First do basic validation (no DB)
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50)),
		validation.Field(&req.Password,
			validation.Required.Error("password is required"),
			validation.Length(8, 0).Error("password must be at least 8 characters long"),
		),
		validation.Field(&req.PhoneNumber,
			validation.Required,
			validation.Match(phoneNumberRegex).Error("phone number is not valid"),
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

	// Check phone number uniqueness with context
	if req.PhoneNumber != "" {
		isUnique, err := v.repo.IsPhoneNumberUnique(ctx, req.PhoneNumber)
		if err != nil {
			return fieldErrors, richerror.New(op).
				WithMessage("failed to check phone number").
				WithKind(richerror.KindUnexpected).
				WithErr(err)
		}
		if !isUnique {
			fieldErrors["phone_number"] = errmsg.ErrorMsgPhoneNumberIsNotUnique
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