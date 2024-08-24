package http

import (
	"errors"
	"log/slog"
	"net/http"
	"tat_hockey_pack/internal/interfaces"
	"tat_hockey_pack/internal/utils"

	"tat_hockey_pack/internal/models"
	"time"
)

type SessionManager struct {
	sessions interfaces.SessionService
	log      *slog.Logger
}

const CookieName = "session_id"

var ErrNoAuth = errors.New("No session found")

func NewSessionManager(sessions interfaces.SessionService, log *slog.Logger) *SessionManager {
	return &SessionManager{
		sessions: sessions,
		log:      log,
	}
}

func (s *SessionManager) Index(w http.ResponseWriter, r *http.Request) {
	_, err := s.Check(r)

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/feeds", http.StatusFound)
}

func (s *SessionManager) Check(r *http.Request) (*models.Session, error) {
	const op = "SessionManager.Check"
	request_id := utils.RequestIDFromContext(r.Context())

	s.log = s.log.With(
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("RequestID", request_id))

	sessionCookie, err := r.Cookie(CookieName)
	s.log.Info("Check session cookie",
		"sessionCookie", sessionCookie.Value)

	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			s.log.Info("CheckSession", "msg", "Session cookie not found")
			return nil, ErrNoAuth
		}
		s.log.Error("CheckSession", "msg", err)
		return nil, err
	}

	if err := sessionCookie.Valid(); err != nil {
		s.log.Error("CheckSession", "sessionCookie.Valid", sessionCookie.Valid(), "msg", err)
		return nil, err
	}

	check, err := s.sessions.Check(r.Context(), sessionCookie.Value)
	if err != nil {
		s.log.Error("CheckSession", "s.sessions.Check", sessionCookie.Value, "msg", err.Error())
		return nil, err
	}

	return check, nil
}

func (s *SessionManager) Create(w http.ResponseWriter, r *http.Request, u interfaces.User) error {

	s.log.Debug("SessionManager.Create", "start", "yes")
	sessionID, err := s.sessions.Create(r.Context(), u)
	if err != nil {
		s.log.Error("Failed to create session", "error", err)
		return err
	}

	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    sessionID,
		Expires:  time.Now().Add(3 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	s.log.Info("Session created and cookie set", "sessionID", sessionID)
	return nil
}

func (s *SessionManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	err := s.sessions.DestroyCurrent(r.Context())
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Удаление куки
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	s.log.Info("DestroySession", "msg", "Session destroyed")
	return nil
}

func (s *SessionManager) DestroyAll(w http.ResponseWriter, r *http.Request, u interfaces.User) error {
	err := s.sessions.DestroyAll(r.Context(), u)
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Удаление куки
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	if err != nil {
		return err
	}
	return nil
}
