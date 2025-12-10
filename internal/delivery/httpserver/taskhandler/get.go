package taskhandler

import (
	"log"
	"net/http"
	"strconv"
	"todoapp/internal/pkg/claim"
	"todoapp/internal/pkg/httpmsgerrorhandler"

	"github.com/labstack/echo/v4"
)

func (h Handler) getTaskByID(c echo.Context) error {

	idParam := c.Param("id")

	taskID, err := strconv.ParseUint(idParam,10,32)

	if err !=nil{
		log.Printf("❌ Invalid task ID: %s", idParam)
		return c.JSON(http.StatusBadRequest,echo.Map{
			"message":"invalid task id",
		})
	}

	claims:=claim.GetClaimsFromEchoContext(c)
	userID:=claims.UserID

	log.Printf("🔵 Getting task ID=%d for user=%d", taskID, userID)

	resp,err:=h.taskSvc.GetTaskByID(c.Request().Context(),uint(taskID),userID)

	if err !=nil{
    log.Printf("❌ Service error: %v", err)
		msg,code:=httpmsgerrorhandler.Error(err)
		return c.JSON(code ,echo.Map{
			"message":msg,
		})
	}
	log.Printf("✅ Task found: %s", resp.Task.Title)
	return  c.JSON(http.StatusOK,echo.Map{
		"message":"task retrived successfully",
		"data":resp,
	})
}