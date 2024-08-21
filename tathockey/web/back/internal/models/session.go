package models

import "time"

type Session struct {
	ID        string    `json:"id" db:"id"`
	UserID    uint64    `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (session *Session) GetSessionID() string {
	return session.ID
}

func (session *Session) GetUserID() uint64 {
	return session.UserID
}
