package logging

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type ContextKey string

var LoggerCtxKey ContextKey = "logger"

func LoggingMiddleware(log *logrus.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newEntry := log.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"time":   time.Now(),
			})

			newEntry.Info("Start logging")

			newEntry.Println("LOLOLO")

			ctx := context.WithValue(r.Context(), LoggerCtxKey, newEntry)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
