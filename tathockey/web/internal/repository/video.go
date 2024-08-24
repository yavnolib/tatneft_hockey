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
	saveVideo = "insert into videos (name) values ($1) returning id;"
	getByID   = `select id, name from videos where id = $1`
	getByName = `select id, name from videos where name = $1`
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
		vi.Name).
		Scan(&id)
	if err != nil {
		v.log.Error("Video.Save", "error", err)
		return 0, err
	}

	return id, nil
}

func (v *Video) GetByID(id uint64) (*models.Video, error) {
	v.log.Debug("Video.GetByID", "hash", id)
	var vi models.Video
	err := v.db.QueryRow(context.Background(), getByID, id).
		Scan(&vi.ID, &vi.Name)
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
		Scan(&vi.ID, &vi.Name)
	if err != nil {
		v.log.Error("Video.GetByName", "error", err.Error())
		return nil, err
	}
	v.log.Debug("Video.GetByName", "video", vi)
	return &vi, nil
}
