package sender

import (
    "context"
    "encoding/json"
    "log/slog"

    "github.com/diyorbek/sentinel/agent/internal/models"
    "github.com/diyorbek/sentinel/agent/internal/producer"
)

type Sender struct {
    producer producer.KafkaProducer
    topic    string
    eventCh  <-chan models.Event
}

func New(producer producer.KafkaProducer, topic string, eventCh <-chan models.Event) *Sender {
    return &Sender{
        producer: producer,
        topic:    topic,
        eventCh:  eventCh,
    }
}

func (s *Sender) Run(ctx context.Context) {
    slog.Info("sender started", "topic", s.topic)

    for {
        select {
        case <-ctx.Done():
            slog.Info("sender stopped")
            return

        case event, ok := <-s.eventCh:
            if !ok {
                return
            }
            s.send(ctx, event)
        }
    }
}

func (s *Sender) send(ctx context.Context, event models.Event) {
    data, err := json.Marshal(event)
    if err != nil {
        slog.Error("marshal event failed", "type", event.Type, "err", err)
        return
    }

    key := []byte(event.AgentID)

    if err := s.producer.Produce(ctx, s.topic, key, data); err != nil {
        slog.Error("kafka produce failed", "type", event.Type, "err", err)
        return
    }

    slog.Debug("event sent", "type", event.Type, "agent", event.AgentID)
}