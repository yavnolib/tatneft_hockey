package repository

import (
	"context"
	"errors"
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
	saveUser               = "insert into users (email, nickname, password) values ($1, $2, $3) returning id;"
	getUserByID            = `select id, nickname, email from users where id = $1;`
	getByNickname          = `select id, nickname,email from users where nickname = $1;`
	getByEmail             = `select id, nickname,email from users where email = $1;`
	updatePass             = `update users set password = $1 where email = $2;`
	checkUserByEmailOrNick = `select id from users where email = $1 or nickname = $2;`

	check = `select id, nickname, password from users where email = $1`
)

func NewUserRepository(db *pgxpool.Pool, log *slog.Logger) *User {
	return &User{
		db:  db,
		log: log,
	}
}

// rowToUser - читаем все данные кроме пароля
func rowToUser(row pgx.Row) (*models.User, error) {
	u := &models.User{}
	err := row.Scan(&u.ID, &u.Nickname, &u.Email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) Save(user *models.User) (int64, error) {
	var id int64
	err := u.db.QueryRow(context.Background(), saveUser, user.Email, user.Nickname, user.Password).
		Scan(&id)
	if err != nil {
		return -1, err
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

func (u *User) CheckUser(email, nickname string) (int64, error) {
	var id int64
	err := u.db.QueryRow(context.Background(), checkUserByEmailOrNick, email, nickname).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return -1, nil
	}
	if err != nil {
		u.log.Error("User.CheckUser", "err", err.Error())
		return -1, err
	}
	u.log.Debug("User.CheckUser", "user", id)

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

func (u *User) UpdatePassword(email string, password []byte) error {
	_, err := u.db.Exec(context.Background(), updatePass, password, email)
	if err != nil {
		u.log.Error("User.UpdatePassword", "err", err.Error())
		return err
	}

	return nil
}

func (u *User) CheckPass(email string) (*models.User, error) {
	m := models.User{
		Email: email,
	}
	err := u.db.QueryRow(context.Background(), check, email).Scan(&m.ID, &m.Nickname, &m.Password)

	if err != nil {
		return nil, err
	}
	return &m, nil
}
