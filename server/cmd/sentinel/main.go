package main

import (
	"log"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/handlers"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/internal/repository/postgres"
	"github.com/diyorbek/sentinel/internal/service"
	"github.com/diyorbek/sentinel/pkg/logger"

	_ "github.com/diyorbek/sentinel/internal/handlers/docs"
)

// @Title Sentinel API
// @Version 1.0
// @Description API server for application
// @host localhost:8080
// @BasePath
func main() {
	cfg := config.Load()
	logger := logger.SetupLog()

	db, err := postgres.DBConnection()
	if err != nil {
		logger.Error(err.Error())
		log.Fatal(err)
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos, cfg, logger)
	handlers := handlers.NewHandler(service, logger)

	router := handlers.InitRoutes(cfg)
	router.Run(":8080")
}