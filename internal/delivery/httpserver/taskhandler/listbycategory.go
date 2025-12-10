package taskhandler

import (
	"log"
	"net/http"
	"strconv"
	"todoapp/internal/pkg/claim"
	"todoapp/internal/pkg/httpmsgerrorhandler"

	"github.com/labstack/echo/v4"
)

func (h Handler) getTasksByCategory(c echo.Context) error {
	log.Printf("🔵 Get tasks by category request received")

	categoryIDParam:=c.Param("category_id")
	categoryID,err:=strconv.ParseUint(categoryIDParam,10,32)

	if err!=nil{
		log.Printf("❌ Invalid category ID: %s", categoryIDParam)
		return c.JSON(http.StatusBadRequest,echo.Map{
			"message":"invalid category id",
		})
	}
	claims:=claim.GetClaimsFromEchoContext(c)
	userID:=claims.UserID

 	log.Printf("🔵 Getting tasks for user: %d, category: %d", userID, categoryID)

	resp,err:=h.taskSvc.GetTasksByCategory(c.Request().Context(),userID,uint(categoryID))

	if err !=nil{
		log.Printf("❌ Service error: %v", err)
		msg,code:=httpmsgerrorhandler.Error(err)
		return c.JSON(code,echo.Map{
			"message":msg,
		})
	}

	log.Printf("✅ Found %d tasks in category", len(resp.Tasks))
	return  c.JSON(http.StatusOK,echo.Map{
		"message":"tasks retived successfully",
		"data":resp,
	})


}