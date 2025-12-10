package taskhandler

import (
	"log"
	"net/http"
	"todoapp/internal/param"
	"todoapp/internal/pkg/claim"
	"todoapp/internal/pkg/httpmsgerrorhandler"

	"github.com/labstack/echo/v4"
)

func (h Handler) createTask(c echo.Context) error {
	var req param.CreateTaskRequest

	log.Printf("🔵 Create task request received")

	if err :=c.Bind(&req);err!=nil{
		log.Printf("❌ Bind error: %v", err)
		return c.JSON(http.StatusBadRequest,echo.Map{
			"message":"Invalid request body",
		})
	}

	claims:=claim.GetClaimsFromEchoContext(c)
	req.UserID=claims.UserID

	log.Printf("🔵 Creating task for user: %d, title: %s", req.UserID, req.Title)

	if fieldErrors,err:=h.taskvalidator.ValidateCreateTask(c.Request().Context(),req);err!=nil{
			log.Printf("❌ Validation error: %v, fields: %v", err, fieldErrors)
		msg, code := httpmsgerrorhandler.Error(err)

		return c.JSON(code,echo.Map{
			"message":msg,
			"errors":fieldErrors,
		})
	}
	log.Printf("🔵 Validation passed, calling service...")

	resp,err:=h.taskSvc.CreateTask(c.Request().Context(),req)

	if err!=nil{
		log.Printf("❌ Service error: %v", err)
		msg,code:=httpmsgerrorhandler.Error(err)
		return  c.JSON(code,echo.Map{
			"message":msg,
		})
	}
	log.Printf("✅ Task created successfully: ID=%d", resp.Task.ID)

	return c.JSON(http.StatusCreated,echo.Map{
		"message":"task created successfully",
		"data":resp,
	})


}