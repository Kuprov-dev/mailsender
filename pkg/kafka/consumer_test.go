package kafka

import (
	"context"
	"simplemailsender/pkg/models"
	"testing"
)

func TestGetMessage(t *testing.T) {
	kp, err := NewProducer(
		[]string{"127.0.0.1:9093"},
		"test-topic",
	)
	if err != nil {
		t.Fatal(err)
	}
	k, err := NewConsumer(
		[]string{"127.0.0.1:9093"},
		"test-topic",
		"test-consumer-group",
	)
	if err != nil {
		t.Fatalf("failed to create a new consumer: %v", err)
		return
	}
	t.Run("commit before processing", func(t *testing.T) {
		if err := kp.SendMessages(
			context.Background(),
			[]*models.Message{
				{
					Key:   []byte("commit-before-processing-test"),       //TODO normal key
					Value: []byte("commit-before-processing-test-value"), //TODO ADD json
				},
			},
			0,
		); err != nil {
			t.Errorf("failed to send messages: %v", err)
			return
		}
		t.Logf("sent messages 1")

		// получение следующего сообщения.
		msg, err := k.GetMessage(context.Background())
		if err != nil {
			t.Errorf("failed to get message: %v", err)
			return
		}

		t.Logf("%+v\n%s\n%s", msg, string(msg.Key), string(msg.Value))

	})

	t.Run("commit after processing", func(t *testing.T) {
		if err := kp.SendMessages(
			context.Background(),
			[]*models.Message{
				{
					Key:   []byte("commit-after-processing-test"),
					Value: []byte("commit-after-processing-test-value"),
				},
			},
			1,
		); err != nil {
			t.Errorf("failed to send messages: %v", err)
			return
		}
		if err != nil {
			t.Errorf("failed to create a new consumer: %v", err)
			return
		}

		// получение следующего сообщения.
		_, err = k.ReadAndCommit(context.Background(), func(msg *models.Message) error {
			t.Logf("%+v\n%s\n%s", msg, string(msg.Key), string(msg.Value))
			return nil
		})
		if err != nil {
			t.Errorf("failed to read and commit message: %v", err)
			return
		}
	})
}
