package userhandler

import (
	"log"
	"net/http"
	"todoapp/internal/param"
	"todoapp/internal/pkg/httpmsgerrorhandler"

	"github.com/labstack/echo/v4"
)

func (h Handler) userLogin(c echo.Context) error {
	var req param.LoginRequest

	// Log incoming request
	log.Printf("🔵 Login request received")

	if err := c.Bind(&req); err != nil {
		log.Printf("❌ Bind error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	log.Printf("🔵 Request bound successfully: phone=%s", req.PhoneNumber)

	if fieldErrors, err := h.userValidator.ValidateLoginRequest(c.Request().Context(), req); err != nil {
		log.Printf("❌ Validation error: %v, fields: %v", err, fieldErrors)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	log.Printf("🔵 Validation passed, calling service...")

	resp, err := h.userSvc.Login(c.Request().Context(), req)
	if err != nil {
		log.Printf("❌ Service error: %v", err)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	log.Printf("✅ Login successful for phone: %s", req.PhoneNumber)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user login successfully",
		"data":    resp,
	})
}
