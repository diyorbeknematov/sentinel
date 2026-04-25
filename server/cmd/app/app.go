package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/diyorbek/sentinel/internal/analyzer"
	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/consumer"
	"github.com/diyorbek/sentinel/internal/handlers"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/internal/repository/postgres"
	"github.com/diyorbek/sentinel/internal/service"
)

type App struct {
	server   *http.Server
	consumer *consumer.Consumer
	logger   *slog.Logger
	analyzer *analyzer.LogAnalyzer
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	db, err := postgres.DBConnection()
	if err != nil {
		return nil, fmt.Errorf("db connection: %w", err)
	}

	repos := repository.NewRepository(db)
	svc := service.NewService(repos, cfg, log)
	handler := handlers.NewHandler(svc, log)

	router := handler.InitRoutes(cfg)

	return &App{
		server: &http.Server{
			Addr:         ":8080",
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		logger: log,
	}, nil
}

func (a *App) RunHTTP() error {
	a.logger.Info("server started", slog.String("addr", a.server.Addr))
	return a.server.ListenAndServe()
}

func (a *App) RunKafka(ctx context.Context) error {
	a.logger.Info("consumer started ...")
	go a.analyzer.StartCleanup(ctx)
	return a.consumer.Run(ctx)
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.consumer.Close(); err != nil {
		a.logger.Warn("consumer close failed", "error", err)
	}
	return a.server.Shutdown(ctx)
}
