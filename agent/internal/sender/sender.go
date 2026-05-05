package sender

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/diyorbek/sentinel/agent/internal/models"
	"github.com/diyorbek/sentinel/agent/internal/producer"
)

type Sender struct {
	producer producer.KafkaProducer
	topic    string
	eventCh  <-chan models.Event

	queue chan models.Event
}

func New(producer producer.KafkaProducer, topic string, eventCh <-chan models.Event) *Sender {
	return &Sender{
		producer: producer,
		topic:    topic,
		eventCh:  eventCh,
		queue:    make(chan models.Event, 1000),
	}
}

func (s *Sender) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case event := <-s.queue:
			s.sendWithRetry(ctx, event)
		}
	}
}

func (s *Sender) Run(ctx context.Context) {
	slog.Info("sender started", "topic", s.topic)

	// worker
	go s.worker(ctx)

	for {
		select {
		case <-ctx.Done():
			slog.Info("sender stopped")
			return

		case event, ok := <-s.eventCh:
			if !ok {
				return
			}
			select {
			case s.queue <- event:
			default:
				slog.Warn("queue full, dropping event", "type", event.Type)
			}
		}
	}
}

func (s *Sender) sendWithRetry(ctx context.Context, event models.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		slog.Error("marshal failed", "type", event.Type, "err", err)
		return
	}

	key := []byte(event.AgentID)

	var lastErr error

	for i := 0; i < 5; i++ {
		lastErr = s.producer.Produce(ctx, s.topic, key, data)
		if lastErr == nil {
			slog.Debug("event sent", "type", event.Type)
			return
		}

		slog.Warn("kafka retry", "try", i+1, "err", lastErr)
		time.Sleep(time.Second * time.Duration(i+1))
	}

	slog.Error("event lost after retries", "type", event.Type, "err", lastErr)
}
