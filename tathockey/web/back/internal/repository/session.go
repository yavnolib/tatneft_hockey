package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"tat_hockey_pack/internal/models"
)

type SessionRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewSessionRepository(db *pgxpool.Pool, log *slog.Logger) *SessionRepository {
	return &SessionRepository{
		db:  db,
		log: log,
	}
}

const (
	createSessionQuery = `insert into sessions (id, user_id) values ($1, $2) returning id`
	destroyBySessionID = ``
	destroyAll         = ``
)

func (s *SessionRepository) Create(session models.Session) (string, error) {
	err := s.db.QueryRow(context.Background(), createSessionQuery, session.ID, session.UserID).
		Scan(&session.ID)
	if err != nil {
		return "", err
	}
	return session.ID, nil
}

func (s *SessionRepository) Destroy(sessionID string) error {
	_, err := s.db.Exec(context.Background(), destroyBySessionID, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionRepository) DestroyAll(userID uint64) error {
	_, err := s.db.Exec(context.Background(), destroyAll, userID)
	if err != nil {
		return err
	}
	return nil
}
