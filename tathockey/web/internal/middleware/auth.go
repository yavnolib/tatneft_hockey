package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"tat_hockey_pack/internal/interfaces"
	"tat_hockey_pack/internal/service/session"
)

var noAuth = map[string]struct{}{
	"/login":  {},
	"/signup": {},
	"/feeds":  {},
	"/post":   {},
	"/":       {},
}

func Auth(sm interfaces.SessionManager, log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := noAuth[r.URL.Path]
		log.Info("AuthMiddleware", "auth route", ok)
		if ok || strings.Contains(r.URL.Path, "static") {
			next.ServeHTTP(w, r)
			return
		}

		sess, err := sm.Check(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		ctx := context.WithValue(r.Context(), session.Key, sess)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
