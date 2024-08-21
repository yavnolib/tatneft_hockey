package interfaces

import (
	"net/http"
	"tat_hockey_pack/internal/models"
)

type User interface {
	GetID() int64
}

type UserRepository interface {
	Save(user *models.User) (int64, error)
	CheckUser(email, nickname string) (int64, error)
	UpdatePassword(email string, password []byte) error
	CheckPass(email string) (*models.User, error)
}

type UserService interface {
	Login(email, password string) (int64, error)
	Register(email, nickname, password string) (int64, error)
	UpdatePassword(email, oldPassword, newPassword string) (int64, error)
}

type UserManager interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	UpdatePassword(w http.ResponseWriter, r *http.Request)
}
