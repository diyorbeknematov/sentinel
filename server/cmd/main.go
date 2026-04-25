package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diyorbek/sentinel/cmd/app"
	"github.com/diyorbek/sentinel/internal/config"
	_ "github.com/diyorbek/sentinel/internal/handlers/docs"
	"github.com/diyorbek/sentinel/pkg/logger"
)

// @Title Sentinel API
// @Version 1.0
// @Description API server for application
// @host localhost:8080
// @BasePath /
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()
	logger := logger.SetupLog()

	a, err := app.New(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	// HTTP — goroutine'da
	go func() {
		if err := a.RunHTTP(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// Consumer — goroutine'da
	go func() {
		if err := a.RunKafka(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Signal kutamiz
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := a.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
