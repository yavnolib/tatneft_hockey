package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
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
	getByEmail    = `select id, nickname,email from users where email = $1;`

	checkUserByEmailOrNick = `select id from users where email = $1 or nickname = $2;`
)

func NewUserRepository(db *pgxpool.Pool, log *slog.Logger) *User {
	return &User{
		db:  db,
		log: log,
	}
}

func rowToUser(row pgx.Row) (*models.User, error) {
	u := &models.User{}
	err := row.Scan(&u.ID, &u.Nickname, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
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
	row := u.db.QueryRow(context.Background(), getByNickname, username)

	user, err := rowToUser(row)
	if err != nil {
		u.log.Error("User.GetByNickname", "err", err.Error())
		return nil, err
	}
	u.log.Debug("User.GetByNickname", "user", user)
	return user, nil
}

func (u *User) GetByEmail(email string) (*models.User, error) {
	row := u.db.QueryRow(context.Background(), getByEmail, email)
	user, err := rowToUser(row)
	if err != nil {
		u.log.Error("User.GetByEmail", "err", err.Error())
		return nil, err
	}
	u.log.Debug("User.GetByEmail", "user", user)
	return user, nil
}

func (u *User) CheckUser(email, nickname string) (uint64, error) {
	var id uint64
	err := u.db.QueryRow(context.Background(), checkUserByEmailOrNick, email, nickname).Scan(&id)

	if err != nil {
		u.log.Error("User.GetByEmail", "err", err.Error())
		return 0, err
	}
	u.log.Debug("User.GetByEmail", "user", id)

	return id, nil
}

func (u *User) GetByID(id int64) (*models.User, error) {
	row := u.db.QueryRow(context.Background(), getUserByID, id)

	user, err := rowToUser(row)

	if err != nil {
		u.log.Error("User.GetByID", "err", err.Error())
		return nil, err
	}
	u.log.Debug("User.GetByID", "user", user)
	return user, nil
}
