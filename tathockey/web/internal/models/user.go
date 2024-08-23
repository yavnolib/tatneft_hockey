package models

// User
// Nickname - публичное имя, доступное для изменения
// Email - используется для входа
type User struct {
	ID       int64  `json:"id" db:"id"`
	Nickname string `json:"nickname" db:"nickname"`
	Email    string `json:"email" db:"email"`
	Password []byte `json:"password" db:"password"`
}

func (u *User) GetID() int64 {
	return u.ID
}
