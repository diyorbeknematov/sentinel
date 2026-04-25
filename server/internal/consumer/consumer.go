package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/segmentio/kafka-go"
)

type EventHandler interface {
	HandleMetric(ctx context.Context, event models.Event) error
	HandleAppLog(ctx context.Context, event models.Event) error
	HandleNginxLog(ctx context.Context, event models.Event) error
}

type Consumer struct {
	reader  *kafka.Reader
	handler EventHandler
	logger  *slog.Logger
}

func New(brokers []string, topic, groupID string, handler EventHandler, logger *slog.Logger) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID, // bir xil groupID = load balancing
			MinBytes:       1,
			MaxBytes:       10 << 20, // 10MB
			CommitInterval: 0,        // manual commit — xabar yo'qolmasin
		}),
		handler: handler,
		logger:  logger,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	slog.Info("consumer started")

	for {
		// context bekor bo'lsa to'xtaydi
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("consumer stopped")
				return nil
			}
			return fmt.Errorf("fetch message: %w", err)
		}

		// qayta ishlaymiz
		if err := c.process(ctx, msg); err != nil {
			c.logger.Warn("process failed, skipping", "err", err, "offset", msg.Offset)
		}

		// muvaffaqiyatli o'qildi — Kafka'ga bildirамiz
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Warn("commit failed", "err", err)
		}
	}
}

func (c *Consumer) process(ctx context.Context, msg kafka.Message) error {
	// JSON parse
	var event models.Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	c.logger.Debug("event received", "type", event.Type, "agent", event.AgentID)

	// type ga qarab yo'naltirамiz
	switch event.Type {
	case models.EventMetric:
		return c.handler.HandleMetric(ctx, event)
	case models.EventAppLog:
		return c.handler.HandleAppLog(ctx, event)
	case models.EventNginxLog:
		return c.handler.HandleNginxLog(ctx, event)
	default:
		c.logger.Warn("unknown event type", "type", event.Type)
		return nil
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
