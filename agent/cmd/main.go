package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/diyorbek/sentinel/agent/cmd/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		slog.Error("init failed", "err", err)
		os.Exit(1)
	}
	defer a.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.RunCollectors(ctx)
	a.RunSender(ctx)

	slog.Info("agent started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down...")
	cancel()
}
