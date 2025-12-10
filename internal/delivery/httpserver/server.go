package httpserver

import (
	"fmt"
	"todoapp/internal/config"
	"todoapp/internal/delivery/httpserver/categoryhandler"
	"todoapp/internal/delivery/httpserver/taskhandler"
	"todoapp/internal/delivery/httpserver/userhandler"
	"todoapp/internal/pkg/validator/categoryvalidator"
	"todoapp/internal/pkg/validator/taskvalidator"
	"todoapp/internal/pkg/validator/uservalidator"
	"todoapp/internal/service/authservice"
	"todoapp/internal/service/categoryservice"
	"todoapp/internal/service/taskservice"
	"todoapp/internal/service/userservice"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config config.Config
	userHandler  userhandler.Handler
	categoryHandler categoryhandler.Handler
	taskHandler     taskhandler.Handler
	Router *echo.Echo
}


func New (
	config config.Config,
	authSvc authservice.Service,
	userSvc userservice.Service,
	userValidator uservalidator.Validator,
	categorySvc categoryservice.Service,
	categoryValidator categoryvalidator.Validator,
	taskSvc taskservice.Service,
	taskValidator taskvalidator.Validator,
) Server {
	return  Server{
		Router:echo.New(),
		config: config,
		userHandler: userhandler.New(authSvc, userSvc, userValidator, config.Auth),
		categoryHandler: categoryhandler.New(config.Auth, authSvc, categorySvc, categoryValidator),
		taskHandler: taskhandler.New(
			config.Auth,
			authSvc,
			taskSvc,
			taskValidator,
		),
	}
}


func (s Server) Serve() {
	s.Router = echo.New()
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())

	s.Router.GET("/health-check", s.healthCheck)


	s.userHandler.SetRoutes(s.Router)
	s.categoryHandler.SetRoutes(s.Router)
	s.taskHandler.SetRoutes(s.Router)

		// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error", err)
	}
}