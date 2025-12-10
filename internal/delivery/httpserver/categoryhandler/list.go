package categoryhandler

import (
	"log"
	"net/http"
	"todoapp/internal/pkg/claim"
	"todoapp/internal/pkg/httpmsgerrorhandler"

	"github.com/labstack/echo/v4"
)

func (h Handler) getAllCategories(c echo.Context) error {
	log.Printf("🔵 Get all categories request received")

	claims := claim.GetClaimsFromEchoContext(c)
	userID := claims.UserID

	log.Printf("🔵 Getting categories for user: %d", userID)

	resp,err:=h.categorySvc.GetAllCategories(c.Request().Context(),userID)

	if err!=nil{
		log.Printf("❌ Service error: %v", err)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}
	log.Printf("✅ Found %d categories", len(resp.Categories))

	return c.JSON(http.StatusOK,echo.Map{
		"message":"categories retrived successfully",
		"data":resp,
	})

}