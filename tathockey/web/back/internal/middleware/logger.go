package middleware

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

func requestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return "-"
	}
	return requestID
}

func LoggerMiddleware(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := requestIDFromContext(r.Context())
		if requestID == "" {
			requestID = uuid.New().String()
		}

		log = log.With(
			slog.String("request_id", requestID),
			slog.String("url", r.URL.String()),
			slog.String("method", r.Method),
			slog.String("remote_address", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)
		next.ServeHTTP(w, r)

		log.Info("HTTP request served",
			slog.String("duration", time.Since(start).String()),
		)
	})
}
