// internal/delivery/httpserver/categoryhandler/delete.go
package categoryhandler

import (
	"log"
	"net/http"
	"strconv"
	"todoapp/internal/pkg/claim"
	"todoapp/internal/pkg/httpmsgerrorhandler"

	"github.com/labstack/echo/v4"
)

func (h Handler) deleteCategory(c echo.Context) error {
	log.Printf("🔵 Delete category request received")

	// Get category ID from URL
	idParam := c.Param("id")
	categoryID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Printf("❌ Invalid category ID: %s", idParam)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid category id",
		})
	}

	// Get authenticated user
	claims := claim.GetClaimsFromEchoContext(c)
	userID := claims.UserID

	log.Printf("🔵 Deleting category ID=%d for user=%d", categoryID, userID)

	// Delete category
	_,err = h.categorySvc.DeleteCategory(c.Request().Context(), uint(categoryID), userID)
	if err != nil {
		log.Printf("❌ Service error: %v", err)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	log.Printf("✅ Category deleted successfully: ID=%d", categoryID)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "category deleted successfully",
	})
}