package interfaces

import "tat_hockey_pack/internal/models"

type User interface {
	GetID() int64
}

type UserRepository interface {
	Save(user *models.User) (int64, error)
	GetByNickname(username string) (*models.User, error)
	GetByID(id int64) (*models.User, error)
}
