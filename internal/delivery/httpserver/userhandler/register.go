package userhandler

import (
	"net/http"
	"todoapp/internal/param"
	"todoapp/internal/pkg/httpmsgerrorhandler"

	"github.com/labstack/echo/v4"
)

func (h Handler) userRegister(c echo.Context) error {
	var req param.RegisterRequest
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	if fieldErrors, err := h.userValidator.ValidateRegisterRequest(c.Request().Context(),req); err != nil {
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

		resp, err := h.userSvc.Register(c.Request().Context(), req)
	if err != nil {
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "user registered successfully",
		"data":    resp,
	})
}