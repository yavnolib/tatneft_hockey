package user

import (
	"database/sql"
	"errors"
	"log/slog"
	"tat_hockey_pack/internal/models"
	"tat_hockey_pack/internal/repository"
	"tat_hockey_pack/internal/repository/repo_errors"
)

// TODO можно реализовать замену пароля
// добавить секретное слово(?) почту

var ErrInvalidPass = errors.New("invalid password")

type Service struct {
	log  *slog.Logger
	repo *repository.User
}

func NewService(logger *slog.Logger, UserRepo *repository.User) *Service {
	return &Service{
		log:  logger,
		repo: UserRepo,
	}
}

func (s *Service) Register(email, nickname, password string) (int64, error) {
	salt, err := s.makeSalt()
	if err != nil {
		s.log.Error("Failed to generate salt", "error", err.Error())
		return 0, err
	}
	hashPass, err := s.hashPassword([]byte(password), salt)
	if err != nil {
		s.log.Error("Failed to hash password", "error", err)
		return -1, err
	}
	id, err := s.repo.CheckUser(email, nickname)
	s.log.Debug("User id", "id", id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.Error("Failed to check user", "error", err)
		return 0, err
	}
	if id != -1 {
		s.log.Debug("User already exists", "id", id)
		s.log.Error("User already exists", "email", email)
		return -2, repo_errors.ErrUserExists
	}
	user := &models.User{
		Nickname: nickname,
		Email:    email,
		Password: hashPass,
	}

	id, err = s.repo.Save(user)
	if err != nil {
		s.log.Error("Failed to create user", "error", err)
		return -3, err
	}

	return id, nil
}

func (s *Service) Login(email, password string) (int64, error) {
	user, err := s.repo.CheckPass(email)
	if err != nil {
		return -1, err
	}
	if user.Nickname == "" {
		return -1, repo_errors.ErrUserNotExists
	}
	if !s.passwordIsValid([]byte(password), user.Password) {
		return -1, ErrInvalidPass
	}
	s.log.Debug("User logged in", "user", user)
	return user.ID, nil
}

func (s *Service) UpdatePassword(email, oldPassword, newPassword string) (int64, error) {
	user, err := s.repo.CheckPass(email)
	if err != nil {
		return -1, err
	}
	if user.Nickname == "" {
		return -1, repo_errors.ErrUserNotExists
	}
	if !s.passwordIsValid([]byte(oldPassword), user.Password) {
		return -1, ErrInvalidPass
	}
	salt, err := s.makeSalt()
	if err != nil {
		s.log.Error("UpdatePassword - Failed to generate salt", "error", err.Error())
		return -1, err
	}
	hashPass, err := s.hashPassword([]byte(newPassword), salt)
	err = s.repo.UpdatePassword(email, hashPass)
	if err != nil {
		return -1, err
	}

	return user.ID, nil
}
