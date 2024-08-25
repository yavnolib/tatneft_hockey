package repo_errors

import "errors"

var (
	ErrNoSession     = errors.New("no session")
	ErrUserNotExists = errors.New("user not exists")
	ErrUserExists    = errors.New("user exists")
)
