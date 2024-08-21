package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"tat_hockey_pack/internal/models"
)

//Todo убрать логгер

type User struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

const (
	saveUser      = "insert into users (id, email, nickname, password) values ($1, $2, $3, $4) returning id;"
	getUserByID   = `select id, nickname, email, password from users where id = $1;`
	getByNickname = `select id, nickname,email, password from users where nickname = $1;`
)

func NewUserRepository(db *pgxpool.Pool, log *slog.Logger) *User {
	return &User{
		db:  db,
		log: log,
	}
}

func (u *User) Save(user *models.User) (uint64, error) {
	var id uint64
	err := u.db.QueryRow(context.Background(), saveUser, user.ID, user.Nickname, user.Email, user.Password).
		Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *User) GetByNickname(username string) (*models.User, error) {
	var user models.User
	err := u.db.QueryRow(context.Background(), getByNickname, username).
		Scan(&user.ID, &user.Nickname, &user.Email, &user.Password)

	if err != nil {
		u.log.Error("User.GetByNickname", "err", err.Error())
		return nil, err
	}
	u.log.Debug("User.GetByNickname", "user", user)
	return &user, nil
}

func (u *User) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := u.db.QueryRow(context.Background(), getUserByID, id).
		Scan(&user.ID, &user.Nickname, &user.Email, &user.Password)
	if err != nil {
		u.log.Error("User.GetByID", "err", err.Error())
		return nil, err
	}
	u.log.Debug("User.GetByID", "user", user)
	return &user, nil
}
