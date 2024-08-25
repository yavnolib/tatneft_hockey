package models

const GifPath = "./app/gifs"

type Gif struct {
	ID         int64  `json:"id" db:"id"`
	Name       string `json:"path" db:"path"`
	EventClass string `json:"class_name" db:"class_name"`
}
