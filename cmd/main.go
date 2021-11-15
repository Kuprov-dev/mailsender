package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"simplemailsender/pkg/conf"
	"simplemailsender/pkg/db"
	"syscall"
	"time"
)

func main() {
	conf := conf.New()
	db.ConnectMongoDB(context.TODO(), conf)

	printDAO := db.NewMongoDBTemp(context.TODO(), db.GetMongoDBConnection())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for range time.Tick(time.Second * 10) {
			printDAO.PrintTemplates()

		}
	}()

	<-stop
	fmt.Println("kek")

}
