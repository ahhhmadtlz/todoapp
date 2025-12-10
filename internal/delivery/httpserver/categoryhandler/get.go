package categoryhandler

import (
	"log"
	"net/http"
	"strconv"
	"todoapp/internal/pkg/claim"
	"todoapp/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) getCategoryByID(c echo.Context) error {
	log.Printf("🔵 Get category by ID request received")

	idParam := c.Param("id")
	categoryID, err := strconv.ParseInt(idParam,10,32)

	if err!=nil{
		log.Printf("❌ Invalid category ID: %s", idParam)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid category id",
		})
	}
	claims:=claim.GetClaimsFromEchoContext(c)
	userID:=claims.UserID

	log.Printf("🔵 Getting category ID=%d for user=%d", categoryID, userID)
		// Get category
	resp, err := h.categorySvc.GetCategoryByID(c.Request().Context(), uint(categoryID), userID)
	if err != nil {
		log.Printf("❌ Service error: %v", err)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	log.Printf("✅ Category found: %s", resp.Category.Name)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "category retrieved successfully",
		"data":    resp,
	})
}