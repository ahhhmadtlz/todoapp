package taskhandler

import (
	"log"
	"net/http"
	"todoapp/internal/pkg/claim"
	"todoapp/internal/pkg/httpmsgerrorhandler"

	"github.com/labstack/echo/v4"
)

func (h Handler) getAllTasks(c echo.Context) error {

	log.Printf("🔵 Get all tasks request received")
	claims := claim.GetClaimsFromEchoContext(c)
	userID:=claims.UserID

	log.Printf("🔵 Getting tasks for user: %d", userID)

	resp,err:=h.taskSvc.GetAllTasks(c.Request().Context(),userID)

	if err !=nil {
		log.Printf("❌ Service error: %v", err)
		msg,code :=httpmsgerrorhandler.Error(err)
		return c.JSON(code ,echo.Map{
			"message":msg,
		})
	}
	log.Printf("✅ Found %d tasks", len(resp.Tasks))
	

		return c.JSON(http.StatusOK, echo.Map{
		"message": "tasks retrieved successfully",
		"data":    resp,
	})
}
