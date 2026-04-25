package app

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/diyorbek/sentinel/agent/internal/clients"
	"github.com/diyorbek/sentinel/agent/internal/collector"
	"github.com/diyorbek/sentinel/agent/internal/config"
	"github.com/diyorbek/sentinel/agent/internal/models"
	"github.com/diyorbek/sentinel/agent/internal/producer"
	"github.com/diyorbek/sentinel/agent/internal/sender"
	"github.com/diyorbek/sentinel/agent/pkg"
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

	// agent_id bo'sh bo'lsa — register qilamiz
	if cfg.AgentID == "" {
		res, err := clients.Register(cfg.ServerURL, models.RegisterRequest{
			Name:      cfg.AgentName,
			IPAddress: pkg.GetLocalIP(),
		})
		if err != nil {
			return nil, fmt.Errorf("register: %w", err)
		}

		cfg.AgentID = res.AgentID.String()
		cfg.APIKey = res.APIKey

		// config.yaml ga yozamiz
		if err := config.Save("config.yaml", cfg); err != nil {
			slog.Warn("save config failed", "err", err)
		}

		slog.Info("agent registered", "agent_id", cfg.AgentID)
	}

    //  Kafka producer
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
