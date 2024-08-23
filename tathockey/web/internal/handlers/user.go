package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"tat_hockey_pack/internal/interfaces"
)

type UserManager struct {
	sessions interfaces.SessionService
	users    interfaces.UserService
	log      *slog.Logger
}

func NewUserManager(
	sessions interfaces.SessionService,
	users interfaces.UserService,
	log *slog.Logger) *UserManager {
	return &UserManager{
		sessions: sessions,
		users:    users,
		log:      log,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Answer int    `json:"answer"`
	Token  string `json:"token"`
}

func (u *UserManager) Login(w http.ResponseWriter, r *http.Request) {
	u.log.Debug(
		"method", "Login",
	)
	if r.Method != "POST" {
		return
	}

	var req LoginRequest
	res := LoginResponse{
		Token:  uuid.New().String(),
		Answer: 200,
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	responseJSON(w, res)
}

func responseJSON(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func (u *UserManager) Logout(w http.ResponseWriter, r *http.Request) {
	u.log.Debug(
		"method", "Logout")
	if r.Method != "POST" {
		return
	}

}

func (u *UserManager) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

}

func (u *UserManager) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

}
