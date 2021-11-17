package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"

	"simplemailsender/pkg/models"
)

type Consumer struct {
	Reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) (*Consumer, error) {
	if len(brokers) == 0 || len(brokers[0]) == 0 || len(topic) == 0 || len(groupID) == 0 {
		return nil, errors.New("не указаны параметры подключения к Kafka")
	}

	return &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 1,
			MaxBytes: 10e6,
			//Partition: 0,
		}),
	}, nil
}

func (c *Consumer) GetMessage(ctx context.Context) (*models.Message, error) {
	kafkaMsg, err := c.Reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	return &models.Message{
		Key:   kafkaMsg.Key,
		Value: kafkaMsg.Value,
	}, err
}

func (c *Consumer) ReadAndCommit(ctx context.Context, fn func(m *models.Message) error) (*models.Message, error) {
	kafkaMsg, err := c.Reader.FetchMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}
	log.Println("offset: ", kafkaMsg.Offset)
	m := &models.Message{
		Key:   kafkaMsg.Key,
		Value: kafkaMsg.Value,
	}
	if fn != nil {
		if err := fn(m); err != nil {
			return nil, fmt.Errorf("failed to process message: %w", err)
		}
	}
	if err := c.Reader.CommitMessages(ctx, kafkaMsg); err != nil {
		return nil, fmt.Errorf("failed to commit the read message: %w", err)
	}
	return m, nil
}
