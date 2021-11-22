package messageQueue

import (
	"context"
	"simplemailsender/pkg/kafka"
	"simplemailsender/pkg/models"
)

type Producer interface {
	SendMessages(ctx context.Context, messages []*models.Message, partition int) error
}

func NewProducer(brokers []string, topic string) (Producer, error) {
	return kafka.NewProducer(brokers, topic)
}

type Consumer interface {
	GetMessage(ctx context.Context) (*models.Message, error)
	ReadAndCommit(ctx context.Context, fn func(m *models.Message) error) (*models.Message, error)
}

func NewConsumer(brokers []string, topic string, groupID string) (Consumer, error) {
	return kafka.NewConsumer(brokers, topic, groupID)
}
