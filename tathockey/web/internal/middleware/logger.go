package middleware

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

type ContextKey string

const (
	RequestIDKey contextKey = "requestID"
)

func RequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return "-"
	}
	return requestID
}

func LoggerMiddleware(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate new request ID if not present
		requestID := RequestIDFromContext(r.Context())
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Create new context with request ID
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Update the logger with request ID
		log = log.With(
			slog.String("request_id", requestID),
			slog.String("url", r.URL.String()),
			slog.String("method", r.Method),
			slog.String("remote_address", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)

		// Use the updated context
		next.ServeHTTP(w, r.WithContext(ctx))

		// Log the request
		log.Info("HTTP request served",
			slog.String("duration", time.Since(start).String()),
		)
	})
}
