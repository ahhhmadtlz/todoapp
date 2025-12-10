// internal/delivery/httpserver/categoryhandler/route.go
package categoryhandler

import (
	"todoapp/internal/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	categoryGroup := e.Group("/categories")
	
	// Apply auth middleware to all category routes
	categoryGroup.Use(middleware.Auth(h.authSvc, h.authConfig))

	categoryGroup.POST("", h.createCategory)
	categoryGroup.GET("", h.getAllCategories)
	categoryGroup.GET("/:id", h.getCategoryByID)
	categoryGroup.PUT("/:id", h.updateCategory)
	categoryGroup.DELETE("/:id", h.deleteCategory)
}