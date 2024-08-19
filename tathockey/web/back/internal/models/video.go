package models

import "time"

const (
	UploadDir = "./uploads"
)

type Video struct {
	ID          uint64    `db:"id" json:"id"`
	CreatorID   uint64    `db:"creator_id" json:"creator_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Hash        string    `db:"hash" json:"hash"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

var _ = UploadDir
