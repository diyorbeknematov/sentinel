package producer

import (
    "context"
    "fmt"
    "time"

    "github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
    Produce(ctx context.Context, topic string, key, value []byte) error
    Close() error
}

type kafkaProducer struct {
    writer *kafka.Writer
}

func New(brokers []string) (KafkaProducer, error) {
    if len(brokers) == 0 {
        return nil, fmt.Errorf("kafka: at least one broker required")
    }
    return &kafkaProducer{
        writer: &kafka.Writer{
            Addr:         kafka.TCP(brokers...),
            Balancer:     &kafka.LeastBytes{},
            WriteTimeout: 10 * time.Second,
            BatchSize:    100,
            BatchTimeout: 10 * time.Millisecond,
            RequiredAcks: kafka.RequireAll,
            Compression:  kafka.Snappy,
            MaxAttempts:  3,
        },
    }, nil
}

func (p *kafkaProducer) Produce(ctx context.Context, topic string, key, value []byte) error {
    err := p.writer.WriteMessages(ctx, kafka.Message{
        Topic: topic,
        Key:   key,
        Value: value,
        Time:  time.Now(),
    })
    if err != nil {
        return fmt.Errorf("kafka produce to %q: %w", topic, err)
    }
    return nil
}

func (p *kafkaProducer) Close() error {
    return p.writer.Close()
}