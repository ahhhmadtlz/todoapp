// internal/delivery/httpserver/categoryhandler/update.go
package categoryhandler

import (
	"log"
	"net/http"
	"strconv"
	"todoapp/internal/param"
	"todoapp/internal/pkg/claim"
	"todoapp/internal/pkg/httpmsgerrorhandler"

	"github.com/labstack/echo/v4"
)

func (h Handler) updateCategory(c echo.Context) error {
	var req param.UpdateCategoryRequest

	log.Printf("🔵 Update category request received")

	// Get category ID from URL
	idParam := c.Param("id")
	categoryID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Printf("❌ Invalid category ID: %s", idParam)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid category id",
		})
	}

	// Bind request
	if err := c.Bind(&req); err != nil {
		log.Printf("❌ Bind error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	// Get authenticated user
	claims := claim.GetClaimsFromEchoContext(c)
	req.UserID = claims.UserID
	req.ID = uint(categoryID)

	log.Printf("🔵 Updating category ID=%d for user=%d", categoryID, req.UserID)

	// Validate
	if fieldErrors, err := h.categoryValidator.ValidateUpdateCategory(c.Request().Context(), uint(categoryID), req); err != nil {
		log.Printf("❌ Validation error: %v, fields: %v", err, fieldErrors)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	log.Printf("🔵 Validation passed, calling service...")

	// Update category
		
	resp, err := h.categorySvc.UpdateCategory(c.Request().Context(), req)
	if err != nil {
		log.Printf("❌ Service error: %v", err)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	log.Printf("✅ Category updated successfully: ID=%d", categoryID)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "category updated successfully",
		"data":    resp,
	})
}