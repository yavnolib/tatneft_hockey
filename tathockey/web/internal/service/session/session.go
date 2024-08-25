package session

import (
	"context"
	"errors"
	"log/slog"
	"tat_hockey_pack/internal/interfaces"
	"tat_hockey_pack/internal/models"
	"tat_hockey_pack/internal/repository"
	"tat_hockey_pack/internal/utils/ses"
	"time"
)

const (
	delay = time.Hour * 24 * 90
)

var (
	invalidSessionErr = errors.New("invalid session")
)

type Service struct {
	repo *repository.Session
	log  *slog.Logger
}

func NewService(log *slog.Logger, sessionRepo *repository.Session) *Service {
	return &Service{
		repo: sessionRepo,
		log:  log,
	}
}

func (s *Service) Check(ctx context.Context, sessionFromCookie string) (*models.Session, error) {
	ses, err := s.repo.GetBySessionID(ctx, sessionFromCookie)

	if err != nil {
		s.log.Error("Session Service", "error", err.Error(),
			"error", err.Error(),
			"ses", ses)
		return nil, err
	}
	if ses.CreatedAt.Add(delay).Before(time.Now()) {
		s.log.Error("ses.CreatedAt.Add(delay).Before(time.Now())", "error", err)
		return nil, invalidSessionErr
	}

	return ses, nil
}

func (s *Service) Create(ctx context.Context, u interfaces.User) (string, error) {
	s.log.Debug("SessionService.Create", "start", time.Now())
	session := models.Session{
		ID:     ses.RandStringRunes(16),
		UserID: u.GetID(),
	}

	err := s.repo.Create(ctx, &session)
	if err != nil {
		return "", err
	}
	return session.ID, nil
}

// DestroyCurrent для logout
func (s *Service) DestroyCurrent(ctx context.Context) error {
	session, err := ses.FromContext(ctx)
	if err != nil {
		s.log.Error("service.session.DestroyCurrent",
			"func", "FromContext",
			"error", err)
		return err
	}

	err = s.repo.Destroy(ctx, session.ID)
	if err != nil {
		s.log.Debug("service.session.Destroy", "error", err)
		return err
	}
	s.log.Debug("service.session.Destroy",
		"session", session,
		"destroy", "done")
	return nil
}

// DestroyAll в случае смены пароля
func (s *Service) DestroyAll(ctx context.Context, user interfaces.User) error {
	s.log.Debug("service.session.DestroyAll", "user", user)
	err := s.repo.DestroyAll(ctx, user.GetID())
	if err != nil {
		s.log.Debug("service.session.DestroyAll", "error", err)
		return err
	}

	return nil
}
