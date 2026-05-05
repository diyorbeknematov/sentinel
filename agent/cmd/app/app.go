package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

func New(apiKey, server, name string) (*App, error) {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	if name != "" {
		cfg.AgentName = name
	} else if cfg.AgentName == "" {
		hostname, _ := os.Hostname()
		cfg.AgentName = fmt.Sprintf("%s-%s", hostname, pkg.GetLocalIP())

	}

	// API Key bo'sh bo'lsa — register qilamiz
	if cfg.APIKey == "" {
		resp, err := clients.Register(server, apiKey, models.RegisterRequest{
			Name:      cfg.AgentName,
			IPAddress: pkg.GetLocalIP(),
		})
		if err != nil {
			return nil, fmt.Errorf("register: %w", err)
		}

		cfg.AgentID = resp.AgentID.String()
		cfg.APIKey = apiKey
		cfg.ServerURL = server
		cfg.KafkaBrokers = resp.KafkaBrokers
		cfg.KafkaTopic = resp.KafkaTopic

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

func (a *App) Heartbeat(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(30 * time.Second) // har 30s
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := clients.SendHeartbeat(a.cfg); err != nil {
					log.Println("heartbeat error:", err)
				}
			case <-ctx.Done(): // Dastur to'xtatilsa, goroutina ham to'xtaydi
				log.Println("Heartbeat stopped")
				return
			}
		}
	}()
}
