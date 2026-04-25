package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "github.com/diyorbek/sentinel/agent/internal/collector"
    "github.com/diyorbek/sentinel/agent/internal/config"
    "github.com/diyorbek/sentinel/agent/internal/models"
    "github.com/diyorbek/sentinel/agent/internal/producer"
    "github.com/diyorbek/sentinel/agent/internal/sender"
)

func main() {
    // Config
    cfg, err := config.Load("config.yaml")
    if err != nil {
        slog.Error("config load failed", "err", err)
        os.Exit(1)
    }

    // Kafka producer
    prod, err := producer.New(cfg.KafkaBrokers)
    if err != nil {
        slog.Error("kafka init failed", "err", err)
        os.Exit(1)
    }
    defer prod.Close()

    // Channel — barcha eventlar shu yerdan o'tadi
    eventCh := make(chan models.Event, 100)

    // Context — Ctrl+C da hammasi to'xtaydi
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Collectorlar
    go collector.StartMetricsCollector(ctx, cfg, eventCh)
    go collector.StartAppLogCollector(ctx, cfg, eventCh)
    go collector.StartNginxLogCollector(ctx, cfg, eventCh)

    // Sender
    s := sender.New(prod, cfg.KafkaTopic, eventCh)
    go s.Run(ctx)

    slog.Info("agent started", "agent_id", cfg.AgentID)

    // Ctrl+C yoki kill signal kutamiz
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    slog.Info("shutting down...")
    cancel()
}