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
		topic   = "kek"
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
	double := make(chan bool)
	wg := &sync.WaitGroup{}
	defer close(responsesMQ)
	defer close(double)
	ctx, cancelCtx := context.WithCancel(context.TODO())
	defer cancelCtx()
	var m sync.Mutex
	for {
		select {
		case <-ctx.Done():
			log.Println("Done!")
			return nil
		case <-time.After(1 * time.Second):
			wg.Add(3)
			runInterruptor(wg, cancelCtx)

			go consume(ctx, wg, consumer, responsesMQ, &m)

			go produce(ctx, wg, producer, responsesMQ, &m)
		}
	}

}

func produce(ctx context.Context, wg *sync.WaitGroup, producer messageQueue.Producer, ch chan models.TemplateMQ, mt *sync.Mutex) error {
	defer wg.Done()
	var update models.TemplateMQ
	request := <-ch
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	msgCount := 0
	for {
		select {
		case <-ticker.C:
			//mt.Lock()
			if request.Status == "pending" {
				request.Status = "done"
				log.Printf("send")

				data, errEnc := json.Marshal(request)
				if errEnc != nil {
					return fmt.Errorf("wron smth %w", errEnc)
				}
				msg := &models.Message{
					Value: data,
				}
				log.Printf("sending message %s", msg.Value)
				if err := producer.SendMessages(ctx, []*models.Message{msg}, 1); err != nil {
					log.Printf("failed to send message #%d: %v", msgCount, err)
					continue
				}

				msgCount++
				ch <- update

			}
			//mt.Unlock()
		case <-ctx.Done():
			return nil
		}
	}
}

func consume(ctx context.Context, wg *sync.WaitGroup, consumer messageQueue.Consumer, ch chan models.TemplateMQ, mt *sync.Mutex) error {
	defer wg.Done()
	var request models.TemplateMQ
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			mt.Lock()
			mes, err := consumer.GetMessage(ctx)
			if err != nil {
				return fmt.Errorf("not valed %w", err)
			}
			errjson := json.Unmarshal(mes.Value, &request)
			if errjson != nil {
				return fmt.Errorf("not valed %w", errjson)
			}
			log.Println(request)
			ch <- request
			mt.Unlock()
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
