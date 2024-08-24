package http

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"tat_hockey_pack/internal/interfaces"
	"tat_hockey_pack/internal/models"
	"tat_hockey_pack/internal/repository/repo_errors"
	"time"
)

type UserManager struct {
	sessions interfaces.SessionManager
	users    interfaces.UserService
	log      *slog.Logger
}

func NewUserManager(
	sessions interfaces.SessionManager,
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
	Status  int    `json:"status"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

func (u *UserManager) Login(w http.ResponseWriter, r *http.Request) {
	u.log.Debug("UserManager.Login", "start", time.Now())
	if r.Method != "POST" {
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.log.Error("Login Error", "err", err.Error(), "req", req)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := u.users.Login(req.Email, req.Password)
	if errors.Is(err, repo_errors.ErrUserNotExists) {
		responseJSON(w, LoginResponse{
			Status:  http.StatusUnauthorized,
			Token:   "",
			Message: "User does not exist",
		})
		return
	}

	if err != nil {
		u.log.Error("Login Handler Error", "err", err.Error())
		responseJSON(w, LoginResponse{
			Status:  http.StatusUnauthorized,
			Token:   "",
			Message: "Failed to login",
		})
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = u.sessions.Create(w, r, &models.User{ID: user})
	if err != nil {
		u.log.Error("Login Handler Error -- Sessions", "err", err.Error())
		return
	}

	u.log.Info("Login Success", "user", user)
	responseJSON(w, LoginResponse{
		Status:  http.StatusOK,
		Token:   "session_token_here",
		Message: "Login succeeded",
	})
}

type RegisterResponse struct {
}

func responseJSON(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	switch body := body.(type) {
	case LoginResponse:
		json.NewEncoder(w).Encode(body)
	case RegisterResponse:
		json.NewEncoder(w).Encode(body)
	}

}

func (u *UserManager) Logout(w http.ResponseWriter, r *http.Request) {
	err := u.sessions.DestroyCurrent(w, r)
	if err != nil {
		u.log.Error("Logout Handler Error -- Sessions", "err", err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

type RegisterRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserManager) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseJSON(w, LoginResponse{
			Status:  http.StatusBadRequest,
			Message: "err json decode",
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u.log.Debug("Register", "request", req)
	user, err := u.users.Register(req.Email, req.Nickname, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		responseJSON(w, LoginResponse{
			Status:  http.StatusUnauthorized,
			Message: "err register",
		})
		return
	}
	err = u.sessions.Create(w, r, &models.User{ID: user})
	w.WriteHeader(http.StatusOK)
	responseJSON(w, LoginResponse{
		Status:  http.StatusOK,
		Token:   uuid.New().String(),
		Message: "Register succeeded",
	})
}

func (u *UserManager) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

}
