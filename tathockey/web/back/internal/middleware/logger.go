package middleware

import (
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

func LoggerMiddleware(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := RequestIDFromContext(r.Context())
		if requestID == "" {
			requestID = uuid.New().String()
		}

		log = log.With(
			slog.String("request_id", requestID),
			slog.String("url", r.URL.String()),
		)
		next.ServeHTTP(w, r)

		log.Info("HTTP request served",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_address", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
			slog.String("duration", time.Since(start).String()),
		)
	})
}
