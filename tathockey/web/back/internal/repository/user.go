package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"tat_hockey_pack/internal/models"
)

type UserRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

const (
	saveUser      = "insert into users (id, nickname, password) values ($1, $2, $3) returning id;"
	getUserByID   = `select id, nickname, password from users where id = $1;`
	getByNickname = `select id, nickname, password from users where nickname = $1;`
)

func NewUserRepository(db *pgxpool.Pool, log *slog.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

func (u *UserRepository) Save(user *models.User) (uint64, error) {
	var id uint64
	err := u.db.QueryRow(context.Background(), saveUser, user.ID, user.Nickname, user.Password).
		Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *UserRepository) GetByNickname(username string) (*models.User, error) {
	var user models.User
	err := u.db.QueryRow(context.Background(), getByNickname, username).
		Scan(&user.ID, &user.Nickname, &user.Password)

	if err != nil {
		u.log.Error("UserRepository.GetByNickname", "err", err.Error())
		return nil, err
	}
	u.log.Debug("UserRepository.GetByNickname", "user", user)
	return &user, nil
}

func (u *UserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := u.db.QueryRow(context.Background(), getUserByID, id).Scan(&user.ID, &user.Nickname, &user.Password)
	if err != nil {
		u.log.Error("UserRepository.GetByID", "err", err.Error())
		return nil, err
	}
	u.log.Debug("UserRepository.GetByID", "user", user)
	return &user, nil
}
