package session

import (
	"context"
	"errors"
	"math/rand"
	"tat_hockey_pack/internal/models"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	ErrNoAuth   = errors.New("no session found")
)

const Key int = 1

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func FromContext(ctx context.Context) (*models.Session, error) {
	sess, ok := ctx.Value(Key).(*models.Session)
	if !ok {
		return nil, ErrNoAuth
	}
	return sess, nil
}
