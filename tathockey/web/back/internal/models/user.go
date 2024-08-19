package models

type User struct {
	ID       uint64 `json:"id" db:"id"`
	Nickname string `json:"nickname" db:"nickname"`
	Password []byte `json:"password" db:"password"`
}
