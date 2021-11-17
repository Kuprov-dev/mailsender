package kafka

import (
	"context"
	"simplemailsender/pkg/models"
	"testing"
)

func TestSendMessages(t *testing.T) {
	// Инициализация клиента Kafka.
	kp, err := NewProducer(
		[]string{"127.0.0.1:9093"},
		"test-topic",
	)
	if err != nil {
		t.Fatal(err)
	}

	// Массив сообщений для отправки.
	messages := []*models.Message{
		{
			Key:   []byte("Test Key"),
			Value: []byte("Test Value"),
		},
	}

	// Отправка сообщения.
	if err = kp.SendMessages(context.Background(), messages, 0); err != nil {
		t.Fatal(err)
	}
	t.Logf("sent messages 3")
}
