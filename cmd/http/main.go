package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"todoapp/internal/config"
	"todoapp/internal/delivery/httpserver"
	"todoapp/internal/pkg/validator/categoryvalidator"
	"todoapp/internal/pkg/validator/taskvalidator"
	"todoapp/internal/pkg/validator/uservalidator"
	"todoapp/internal/repository/migrator"
	"todoapp/internal/repository/mysql"
	"todoapp/internal/repository/mysql/mysqlcategory"
	"todoapp/internal/repository/mysql/mysqltask"
	"todoapp/internal/repository/mysql/mysqluser"
	"todoapp/internal/service/authservice"
	"todoapp/internal/service/categoryservice"
	"todoapp/internal/service/taskservice"
	"todoapp/internal/service/userservice"
)

func main() {
	// Load config
	cfg := config.Load("config.yml")
	fmt.Printf("cfg: %+v\n", cfg)

	// Setup database and run migrations
	mysqlDB := mysql.New(cfg.Mysql)
	log.Println("✓ Database connected")

	if err := runMigrations(mysqlDB, cfg); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	
	log.Println("✓ Migrations applied successfully")

	// Setup services
	authSvc, userSvc, userValidator, categorySvc, categoryValidator, taskSvc, taskValidator := setupServices(cfg, mysqlDB)

	// Start HTTP server
	server := httpserver.New(cfg, authSvc, userSvc, userValidator, categorySvc, categoryValidator, taskSvc, taskValidator)

	go func() {
		log.Printf("🚀 Server starting on port %d...\n", cfg.HTTPServer.Port)
		server.Serve()
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("🛑 Received interrupt signal, shutting down gracefully...")

	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cfg.Application.GracefulShutdownTimeout)
	defer cancel()

	if err := server.Router.Shutdown(ctxWithTimeout); err != nil {
		log.Printf("❌ HTTP server shutdown error: %v\n", err)
	} else {
		log.Println("✓ Server shut down successfully")
	}

	<-ctxWithTimeout.Done()
}

func runMigrations(mysqlDB *mysql.MySQLDB, cfg config.Config) error {
	mgr := migrator.New(mysqlDB, migrator.Config{
		MigrationsDir:  "./internal/repository/mysql/migrations",
		MigrationTable: "schema_migrations",
	})

	return mgr.Up(0)
}

func setupServices(cfg config.Config, mysqlDB *mysql.MySQLDB) (
	authservice.Service,
	userservice.Service,
	uservalidator.Validator,
	categoryservice.Service,
	categoryvalidator.Validator,
	taskservice.Service,
	taskvalidator.Validator,
) {
	//auth service
	authSvc := authservice.New(cfg.Auth)
	//user service
	userMysql := mysqluser.New(mysqlDB)
	userSvc := userservice.New(authSvc, userMysql)
	userValidator := uservalidator.New(userMysql)
	//category service
	categoryMysql:=mysqlcategory.New(mysqlDB)
	categorySvc:=categoryservice.New(categoryMysql)
	categoryValidator:=categoryvalidator.New(categoryMysql)

	taskMysql:=mysqltask.New(mysqlDB)
	taskSvc := taskservice.New(taskMysql, categoryMysql)
	taskValidator := taskvalidator.New(taskMysql,categoryMysql) 

	return authSvc, userSvc, userValidator, categorySvc, categoryValidator, taskSvc, taskValidator

}