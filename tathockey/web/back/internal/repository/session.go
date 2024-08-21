package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"tat_hockey_pack/internal/models"
	"tat_hockey_pack/internal/repository/repo_errors"
)

//Todo убрать логгер

type Session struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewSessionRepository(db *pgxpool.Pool, log *slog.Logger) *Session {
	return &Session{
		db:  db,
		log: log,
	}
}

const (
	createSessionQuery = `insert into sessions (id, user_id) values ($1, $2) returning id`
	destroyBySessionID = `delete from sessions where id = $1`
	destroyAll         = `delete from sessions where user_id = $1`
	getBySessionID     = `select * from sessions where id = $1`
)

func (s *Session) GetBySessionID(ctx context.Context, id string) (*models.Session, error) {
	var session models.Session
	err := s.db.QueryRow(ctx, getBySessionID, id).
		Scan(&session.ID, &session.UserID, &session.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo_errors.ErrNoSession
		}
		return nil, err
	}
	return &session, nil
}

func (s *Session) Create(ctx context.Context, session models.Session) (string, error) {
	err := s.db.QueryRow(ctx, createSessionQuery, session.ID, session.UserID).
		Scan(&session.ID)
	if err != nil {
		return "", err
	}

	return session.ID, nil
}

func (s *Session) Destroy(ctx context.Context, sessionID string) error {
	_, err := s.db.Exec(ctx, destroyBySessionID, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) DestroyAll(ctx context.Context, userID uint64) error {
	_, err := s.db.Exec(ctx, destroyAll, userID)
	if err != nil {
		return err
	}
	return nil
}
