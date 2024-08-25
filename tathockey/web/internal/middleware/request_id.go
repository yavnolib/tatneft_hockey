package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const requestIDKey = "requestID"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			r.Header.Set("X-Request-ID", requestID)
			w.Header().Set("X-Request-ID", requestID)
		}
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
