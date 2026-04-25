package app

import (
    "context"
    "fmt"

    "github.com/diyorbek/sentinel/agent/internal/collector"
    "github.com/diyorbek/sentinel/agent/internal/config"
    "github.com/diyorbek/sentinel/agent/internal/models"
    "github.com/diyorbek/sentinel/agent/internal/producer"
    "github.com/diyorbek/sentinel/agent/internal/sender"
    "log/slog"
)

type App struct {
    cfg      *config.Config
    producer producer.KafkaProducer
    sender   *sender.Sender
    eventCh  chan models.Event
}

func New() (*App, error) {
    cfg, err := config.Load("config.yaml")
    if err != nil {
        return nil, fmt.Errorf("config: %w", err)
    }

    prod, err := producer.New(cfg.KafkaBrokers)
    if err != nil {
        return nil, fmt.Errorf("kafka: %w", err)
    }

    eventCh := make(chan models.Event, 100)

    return &App{
        cfg:      cfg,
        producer: prod,
        sender:   sender.New(prod, cfg.KafkaTopic, eventCh),
        eventCh:  eventCh,
    }, nil
}

func (a *App) RunCollectors(ctx context.Context) {
    slog.Info("collectors started")
    go collector.StartMetricsCollector(ctx, a.cfg, a.eventCh)
    go collector.StartAppLogCollector(ctx, a.cfg, a.eventCh)
    go collector.StartNginxLogCollector(ctx, a.cfg, a.eventCh)
}

func (a *App) RunSender(ctx context.Context) {
    slog.Info("sender started")
    go a.sender.Run(ctx)
}

func (a *App) Close() {
    if err := a.producer.Close(); err != nil {
        slog.Warn("producer close failed", "err", err)
    }
}