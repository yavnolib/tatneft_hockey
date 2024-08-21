package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"tat_hockey_pack/internal/models"
)

//Todo убрать логгер

type Video struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

const (
	saveVideo = "insert into video (creator_id, name, hash, created_at, description) values ($1, $2, $3, $4, $5) returning id;"
	getByHash = `select id, creator_id, name, hash, created_at, description from video where hash = $1`
	getByID   = `select id, creator_id, name, hash, created_at, description from video where id = $1`
	getByName = `select id, creator_id, name, hash, created_at, description from video where name = $1`
)

func NewVideoRepository(db *pgxpool.Pool, logger *slog.Logger) *Video {
	return &Video{
		db:  db,
		log: logger,
	}
}

func (v *Video) Save(vi *models.Video) (uint64, error) {
	v.log.Debug("Video.Save", "--", "run")
	var id uint64
	err := v.db.QueryRow(context.Background(), saveVideo,
		vi.CreatorID, vi.Name, vi.Hash, vi.CreatedAt, vi.Description).
		Scan(&id)
	if err != nil {
		v.log.Error("Video.Save", "error", err)
		return 0, err
	}

	return id, nil
}

func (v *Video) GetByHash(hash string) (*models.Video, error) {
	v.log.Debug("Video.GetByHash", "hash", hash)
	var vi models.Video
	err := v.db.QueryRow(context.Background(), getByHash, hash).
		Scan(&vi.ID, &vi.CreatorID, &vi.Name, &vi.Hash, &vi.CreatedAt, &vi.Description)
	if err != nil {
		v.log.Error("Video.GetByHash", "error", err.Error())
		return nil, err
	}
	v.log.Debug("Video.GetByHash", "video", vi)
	return &vi, nil
}

func (v *Video) GetByID(id uint64) (*models.Video, error) {
	v.log.Debug("Video.GetByID", "hash", id)
	var vi models.Video
	err := v.db.QueryRow(context.Background(), getByID, id).
		Scan(&vi.ID, &vi.CreatorID, &vi.Name, &vi.Hash, &vi.CreatedAt, &vi.Description)
	if err != nil {
		v.log.Error("Video.GetByID", "error", err.Error())
		return nil, err
	}
	v.log.Debug("Video.GetByID", "video", vi)
	return &vi, nil
}

func (v *Video) GetByName(name string) (*models.Video, error) {
	v.log.Debug("Video.GetByName", "hash", name)
	var vi models.Video
	err := v.db.QueryRow(context.Background(), getByName, name).
		Scan(&vi.ID, &vi.CreatorID, &vi.Name, &vi.Hash, &vi.CreatedAt, &vi.Description)
	if err != nil {
		v.log.Error("Video.GetByName", "error", err.Error())
		return nil, err
	}
	v.log.Debug("Video.GetByName", "video", vi)
	return &vi, nil
}
