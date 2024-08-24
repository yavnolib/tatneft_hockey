package models

const (
	UploadDir = "./uploads"
)

type Video struct {
	ID   uint64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

var _ = UploadDir
