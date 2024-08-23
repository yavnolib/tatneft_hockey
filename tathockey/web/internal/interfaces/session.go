package interfaces

import (
	"context"
	"net/http"
	"tat_hockey_pack/internal/models"
)

type Session interface {
	GetSessionID() string
	GetUserID() int64
}

type SessionRepository interface {
	GetBySessionID(ctx context.Context, id string) (*models.Session, error)
	Create(ctx context.Context, session *models.Session) error
	Destroy(ctx context.Context, sessionID string) error
	DestroyAll(ctx context.Context, userID uint64) error
}

type SessionService interface {
	Check(ctx context.Context, sessionFromCookie string) (*models.Session, error)
	Create(ctx context.Context, u User) (string, error)
	DestroyCurrent(ctx context.Context) error
	DestroyAll(ctx context.Context, user User) error
}

type SessionManager interface {
	Check(r *http.Request) (*models.Session, error)
}
