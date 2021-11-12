package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"simplemailsender/pkg/db"
	"syscall"
	"time"
)

func main() {

	db.ConnectMongoDB(context.TODO())

	userDAO := db.NewMongoDBUserDAO(context.TODO(), db.GetMongoDBConnection())
	fmt.Println(userDAO)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for range time.Tick(time.Second * 10) {
			userDAO.PrintTemplates()

		}
	}()

	<-stop
	fmt.Println("kek")

}
