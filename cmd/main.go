package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"simplemailsender/pkg/db"
	"syscall"
)

func main() {

	db.ConnectMongoDB(context.TODO())
	//log := logrus.New()
	//log.SetFormatter(&logrus.JSONFormatter{})
	//logEntry := logrus.NewEntry(log)
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
		for {
			userDAO.PrintTemplates()

		}
	}()
	<-stop
	fmt.Println("kek")

	//handler := logging.LoggingMiddleware(logEntry)(r)
	//s := &http.Server{
	//	Addr:    ":8081",
	//	Handler: handler,
	//}
	//defer s.Close()
	//go func() {
	//	fmt.Println("Starting server")
	//	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		log.Println(err)
	//		return
	//	}
	//}()
	//
	//<-stop

}
