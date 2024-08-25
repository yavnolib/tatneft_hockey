package post

import (
	"log/slog"
	"tat_hockey_pack/internal/repository"
)

type Service struct {
	videoRepo repository.Video
	postRepo  repository.Post
	log       slog.Logger
}
