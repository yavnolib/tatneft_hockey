package models

import "time"

const PreviewUpload = "./preview"

type Post struct {
	ID        int64     `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	GIFs      []Gif     `json:"gifs" db:"gifs"`
	Preview   string    `json:"preview" db:"preview"`
	VideoID   int64     `json:"video_id" db:"video_id"`
	CreatorID int64     `json:"creator_id" db:"creator_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type PostPreview struct {
	ID        int64     `json:"id"`
	Preview   string    `json:"preview"`
	CreatorID int64     `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
}
