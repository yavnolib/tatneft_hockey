package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	http2 "tat_hockey_pack/internal/handlers/http_handlers"
	"tat_hockey_pack/internal/interfaces"
	"tat_hockey_pack/internal/utils/ses"
)

var NoAuth = map[string]struct{}{
	"/login":         {},
	"/api/v1/logout": {},
	"/api/v1/login":  {},
	"/api/v1/signup": {},
	"/signup":        {},
	"/feeds":         {},
	"/post":          {},
	"/":              {},
}

type contextKey string

const CookieKey contextKey = "user-cookie"

func Auth(sm interfaces.SessionManager, log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверка на необходимость аутентификации
		if _, ok := NoAuth[r.URL.Path]; ok || strings.Contains(r.URL.Path, "static") {
			log.Info("AuthMiddleware", "auth route", true)
			next.ServeHTTP(w, r)
			return
		}

		// Получаем куку
		cookie, err := r.Cookie(http2.CookieName)
		if cookie == nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.Error("AuthMiddleware", "cookie error", err.Error())

			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		// Проверка сессии
		sess, err := sm.Check(r)
		if err != nil {
			log.Error("AuthMiddleware",
				"session check error", err.Error(),
				"func", "sm.Check")
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		// Сохраняем сессию и куки в контексте
		ctx := context.WithValue(r.Context(), ses.Key, sess)
		ctx = context.WithValue(ctx, CookieKey, cookie)
		r = r.WithContext(ctx)

		log.Debug("AuthMiddleware", "cookie", cookie.Value)

		next.ServeHTTP(w, r)
	})
}
