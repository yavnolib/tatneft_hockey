package user

import (
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/scrypt"
)

func (s *Service) passwordIsValid(inputPass, dbPass []byte) bool {
	salt, err := s.getSalt(dbPass)
	if err != nil {
		return false
	}
	inputHash, err := s.hashPassword(inputPass, salt)
	if err != nil {
		return false
	}
	return string(inputHash) == string(dbPass)
}

func (s *Service) hashPassword(password, salt []byte) ([]byte, error) {
	hash, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
	if err != nil {
		s.log.Error("UserService.utils.hashPass", "error", err.Error())
		return nil, err
	}
	return append(salt, hash...), nil
}

var invalidPass = errors.New("invalid hash password")

func (s *Service) getSalt(hashPass []byte) ([]byte, error) {
	if len(hashPass) < 40 {
		return nil, invalidPass
	}

	salt := hashPass[:8]
	return salt, nil
}

func (s *Service) makeSalt() ([]byte, error) {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)

	if err != nil {
		s.log.Error("UserService.utils.makeSalt",
			"error", err.Error())
		return nil, err
	}
	return salt, nil
}
