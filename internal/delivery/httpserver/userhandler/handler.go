package userhandler

import (
	"todoapp/internal/pkg/validator/uservalidator"
	"todoapp/internal/service/authservice"
	"todoapp/internal/service/userservice"
)

type Handler struct {
	authSvc authservice.Service
	userSvc userservice.Service
	userValidator uservalidator.Validator
	authConfig authservice.Config
}


func New(
	authSvc authservice.Service,
	userSvc userservice.Service,
	userValidator uservalidator.Validator,
	authConfig authservice.Config) Handler{
		return Handler{
			authSvc: authSvc,
			userSvc: userSvc,
			userValidator: userValidator,
			authConfig: authConfig ,
		}
	}

