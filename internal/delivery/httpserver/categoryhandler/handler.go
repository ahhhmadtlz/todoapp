package categoryhandler

import (
	"todoapp/internal/pkg/validator/categoryvalidator"
	"todoapp/internal/service/authservice"
	"todoapp/internal/service/categoryservice"
)

type Handler struct {
	authConfig authservice.Config
	authSvc authservice.Service
	categorySvc categoryservice.Service
	categoryValidator categoryvalidator.Validator
}

func New(
	authConfig authservice.Config,
	authSvc authservice.Service,
	categorySvc categoryservice.Service,
	categoryValidator categoryvalidator.Validator,
) Handler {
	return Handler{
		authConfig:        authConfig,
		authSvc:           authSvc,
		categorySvc:       categorySvc,
		categoryValidator: categoryValidator,
	}
}