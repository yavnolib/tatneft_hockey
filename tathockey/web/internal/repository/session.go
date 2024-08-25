package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"tat_hockey_pack/internal/models"
	"tat_hockey_pack/internal/repository/repo_errors"
	"tat_hockey_pack/internal/utils/ses"
	"time"
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
	getBySessionID     = `select id, user_id, created_at from sessions where id = $1`

	getUserIDbySessionID = `select user_id from sessions where id = $1`
)

func (s *Session) GetUserIDbySessionID(ctx context.Context) (int64, error) {
	var id int64
	sess, _ := ses.FromContext(ctx)
	err := s.db.QueryRow(ctx, getUserIDbySessionID, sess.ID).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *Session) GetBySessionID(ctx context.Context, id string) (*models.Session, error) {
	var sess models.Session
	err := s.db.QueryRow(ctx, getBySessionID, id).
		Scan(&sess.ID, &sess.UserID, &sess.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo_errors.ErrNoSession
		}

		return nil, err
	}
	return &sess, nil
}

func (s *Session) Create(ctx context.Context, session *models.Session) error {
	s.log.Info("SessionRepo.Create",
		"start", time.Now(),
		"user_id", session.UserID,
		"session_id", session.ID)
	err := s.db.QueryRow(ctx, createSessionQuery, session.ID, session.UserID).Scan(&session.ID)
	if err != nil {
		s.log.Error("Failed to create session", "error", err.Error())
		return err
	}
	return nil
}

func (s *Session) Destroy(ctx context.Context, sessionID string) error {
	_, err := s.db.Exec(ctx, destroyBySessionID, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) DestroyAll(ctx context.Context, userID int64) error {
	_, err := s.db.Exec(ctx, destroyAll, userID)
	if err != nil {
		return err
	}
	return nil
}
