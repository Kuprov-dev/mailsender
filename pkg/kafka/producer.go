package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"

	"simplemailsender/pkg/models"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	if len(brokers) == 0 || brokers[0] == "" || topic == "" {
		return nil, errors.New("не указаны параметры подключения к Kafka")
	}

	return &Producer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers[0]),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}, nil
}

func (p *Producer) SendMessages(ctx context.Context, messages []*models.Message, partition int) error {
	for _, m := range messages {
		if err := p.Writer.WriteMessages(ctx, kafka.Message{
			Key:       m.Key,
			Value:     m.Value,
			Partition: partition,
		}); err != nil {
			return fmt.Errorf("failed to send messages to Kafka: %w", err)
		}
	}
	return nil
}
