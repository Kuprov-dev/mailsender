package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"encoding/json"
	"simplemailsender/pkg/messageQueue"
	"simplemailsender/pkg/models"
	"sync"

	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {
	brokers := []string{"127.0.0.1:9093"}
	const (
		topic   = "my-topic"
		groupID = "my-group"
	)
	producer, err := messageQueue.NewProducer(brokers, topic)
	if err != nil {
		return fmt.Errorf("failed to create a new producer: %w", err)
	}
	consumer, err := messageQueue.NewConsumer(brokers, topic, groupID)
	if err != nil {
		return fmt.Errorf("failed to create a new consumer: %w", err)
	}
	responsesMQ := make(chan models.TemplateMQ)
	wg := &sync.WaitGroup{}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	wg.Add(1)
	runInterruptor(wg, cancelCtx)
	wg.Add(1)
	go produce(ctx, wg, producer, responsesMQ)
	wg.Add(1)
	go consume(ctx, wg, consumer, responsesMQ)

	wg.Wait()
	log.Println("Done!")
	return nil
}

func produce(ctx context.Context, wg *sync.WaitGroup, producer messageQueue.Producer, ch chan models.TemplateMQ) error {
	defer wg.Done()
	request := <-ch
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	msgCount := 0
	for {
		select {
		case <-ticker.C:
			request.Status = "done"
			data, errEnc := json.Marshal(request)
			if errEnc != nil {
				return fmt.Errorf("wron smth %w", errEnc)
			}
			msg := &models.Message{
				Value: data,
			}
			log.Printf("sending message %s", msg.Value)
			if err := producer.SendMessages(ctx, []*models.Message{msg}, 0); err != nil {
				log.Printf("failed to send message #%d: %v", msgCount, err)
				continue
			}
			msgCount++
		case <-ctx.Done():
			return nil
		}
	}
}

func consume(ctx context.Context, wg *sync.WaitGroup, consumer messageQueue.Consumer, ch chan models.TemplateMQ) error {
	defer wg.Done()
	var request models.TemplateMQ
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			mes, err := consumer.ReadAndCommit(ctx, func(m *models.Message) error {
				log.Printf("message received: %s: %s", string(m.Key), string(m.Value))

				return nil
			})
			if err != nil {
				log.Printf("failed to get new message: %v", err)
				continue

			}
			errjson := json.Unmarshal(mes.Value, &request)
			if errjson != nil {
				return fmt.Errorf("not valed %w", errjson)
			}
			ch <- request

		case <-ctx.Done():
			return nil
		}
	}
}

func runInterruptor(wg *sync.WaitGroup, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		defer wg.Done()
		<-c
		log.Println("got interruption signal")
		cancel()
	}()
}
