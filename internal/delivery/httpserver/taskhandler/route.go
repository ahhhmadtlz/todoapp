package taskhandler

import (
	"todoapp/internal/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	taskGroup := e.Group("/tasks")

	taskGroup.Use(middleware.Auth(h.authSvc,h.authConfig))

	taskGroup.POST("",h.createTask)
	taskGroup.GET("",h.getAllTasks)
	taskGroup.GET("/:id",h.getTaskByID)
	taskGroup.GET("/category/:category_id",h.getTasksByCategory)
	taskGroup.PUT("/:id",h.updateTask)
	taskGroup.DELETE("/:id",h.deleteTask)

}