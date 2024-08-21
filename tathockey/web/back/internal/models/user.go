package models

// User
// Nickname - публичное имя, доступное для изменения
// Email - используется для входа
type User struct {
	ID       uint64 `json:"id" db:"id"`
	Nickname string `json:"nickname" db:"nickname"`
	Email    string `json:"email" db:"email"`
	Password []byte `json:"password" db:"password"`
}

func (u *User) GetID() uint64 {
	return u.ID
}
