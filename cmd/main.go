package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"simplemailsender/pkg/db"
	"simplemailsender/pkg/mailsender"

	"net/http"
	//"mailsender/pkg/database"
	//"mailsender/pkg/sender"
	//"auth_service/pkg/db"
	logging "simplemailsender/pkg/logging"

	"os"
	"os/signal"
	"syscall"
)

func main() {

	db.ConnectMongoDB(context.TODO())
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	logEntry := logrus.NewEntry(log)
	userDAO := db.NewMongoDBUserDAO(context.TODO(), db.GetMongoDBConnection())
	fmt.Println(userDAO)
	r := mux.NewRouter()
	r.Handle("/templates", mailsender.SenderHandler(userDAO)).Methods(http.MethodGet)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	handler := logging.LoggingMiddleware(logEntry)(r)
	s := &http.Server{
		Addr:    ":8081",
		Handler: handler,
	}
	defer s.Close()
	go func() {
		fmt.Println("Starting server")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
			return
		}
	}()

	<-stop

}
