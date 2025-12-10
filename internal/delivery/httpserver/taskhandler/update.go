package taskhandler

import (
	"log"
	"net/http"
	"strconv"
	"todoapp/internal/param"
	"todoapp/internal/pkg/claim"
	"todoapp/internal/pkg/httpmsgerrorhandler"

	"github.com/labstack/echo/v4"
)

func (h Handler) updateTask(c echo.Context) error {
	var req param.UpdateTaskRequest

	log.Printf("🔵 Update task request received")

	idParam:=c.Param("id")
	taskID,err:=strconv.ParseUint(idParam,10,32)

	if err!=nil{
		log.Printf("❌ Invalid task ID: %s", idParam)
		return c.JSON(http.StatusBadRequest,echo.Map{
			"message":"invalid task id",
		})
	}

	if err:=c.Bind(&req);err !=nil{
		log.Printf("❌ Bind error: %v", err)
		return  c.JSON(http.StatusBadRequest,echo.Map{
			"message":"invalid request body",
		})
	}

	claims:=claim.GetClaimsFromEchoContext(c)
	req.UserID=claims.UserID
	req.ID = uint(taskID)

	log.Printf("🔵 Updating task ID=%d for user=%d", req.ID, req.UserID)

	if fieldErrors, err := h.taskvalidator.ValidateUpdateTask(
		c.Request().Context(),
		req.ID,
		req,
	); err != nil {

		log.Printf("❌ Validation error: %v, fields: %v", err, fieldErrors)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	log.Printf("🔵 Validation passed, calling service...")

	resp, err := h.taskSvc.UpdateTask(c.Request().Context(), req)
	if err != nil {
		log.Printf("❌ Service error: %v", err)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	log.Printf("✅ Task updated successfully: ID=%d", req.ID)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "task updated successfully",
		"data":    resp,
	})
	
}