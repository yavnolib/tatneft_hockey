package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"tat_hockey_pack/internal/interfaces"
	"tat_hockey_pack/internal/models"
	"time"
)

type SessionManager struct {
	sessions interfaces.SessionService
	log      *slog.Logger
}

const cookieName = "session_id"

var ErrNoAuth = errors.New("No session found")

func NewSessionManager(sessions interfaces.SessionService, log *slog.Logger) *SessionManager {
	return &SessionManager{
		sessions: sessions,
		log:      log,
	}
}

func (s *SessionManager) Check(r *http.Request) (*models.Session, error) {
	sessionCookie, err := r.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			s.log.Info("CheckSession", "msg", "Session cookie not found")
			return nil, ErrNoAuth
		}
		return nil, err
	}

	check, err := s.sessions.Check(r.Context(), sessionCookie.Value)
	if err != nil {
		return nil, err
	}

	return check, nil
}

func (s *SessionManager) Create(w http.ResponseWriter, r *http.Request, u interfaces.User) error {
	sessionID, err := s.sessions.Create(r.Context(), u)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    sessionID,
		Expires:  time.Now().Add(3 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	s.log.Debug("CreateSession",
		"sessionID", sessionID,
		"created", "success")
	return nil
}

func (s *SessionManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	err := s.sessions.DestroyCurrent(r.Context())
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     cookieName,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	s.log.Debug("DestroySession",
		"set cookie", cookie,
	)
	return nil
}

func (s *SessionManager) DestroyAll(w http.ResponseWriter, r *http.Request, u interfaces.User) error {
	err := s.sessions.DestroyAll(r.Context(), u)
	if err != nil {
		return err
	}
	return nil
}
