package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"tat_hockey_pack/internal/interfaces"
	"tat_hockey_pack/internal/models"
	"tat_hockey_pack/internal/repository/repo_errors"
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
	u.log.Info("method", "Login")
	u.log.Debug(
		"method", "Login",
	)
	if r.Method != "POST" {
		return
	}

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.log.Error("Login Error", "err", err.Error(),
			"body", r.Body,
			"req", req)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u.log.Error("Login Error",
		"body", r.Body,
		"req", req)
	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u.log.Debug("Login", "request", req)
	user, err := u.users.Login(req.Email, req.Password)
	if errors.Is(err, repo_errors.ErrUserNotExists) {
		responseJSON(w, LoginResponse{
			Status:  http.StatusUnauthorized,
			Token:   "",
			Message: "User does not exist",
		})
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
	//responseJSON(w, LoginResponse{
	//	Status:  http.StatusOK,
	//	Token:   uuid.New().String(),
	//	Message: "Login succeeded",
	//})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(LoginResponse{
		Status:  http.StatusOK,
		Token:   uuid.New().String(),
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
	u.log.Debug(
		"method", "Logout")
	if r.Method != "POST" {
		return
	}

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
	_, err := u.users.Register(req.Email, req.Nickname, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		responseJSON(w, LoginResponse{
			Status:  http.StatusUnauthorized,
			Message: "err register",
		})
		return
	}
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
