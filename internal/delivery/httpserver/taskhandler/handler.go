package taskhandler

import (
	"todoapp/internal/pkg/validator/taskvalidator"
	"todoapp/internal/service/authservice"
	"todoapp/internal/service/taskservice"
)

type Handler struct {
	authConfig authservice.Config
	authSvc authservice.Service
	taskSvc taskservice.Service
	taskvalidator taskvalidator.Validator
}


func New (
	authConfig authservice.Config,
	authSvc authservice.Service,
	taskSvc taskservice.Service,
	taskValidator taskvalidator.Validator,
)Handler {
	return Handler{
		authConfig: authConfig,
		authSvc: authSvc,
		taskSvc: taskSvc,
		taskvalidator: taskValidator,
	}
}