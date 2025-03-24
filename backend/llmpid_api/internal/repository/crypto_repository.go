package repository

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/sirupsen/logrus"
)

type CryptoRepository struct {
	logger *logrus.Logger
}

func NewCryptoRepository(logger *logrus.Logger) *CryptoRepository {
	return &CryptoRepository{logger: logger}
}

func (r *CryptoRepository) GenrateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		r.logger.Error("Unable to generate cryptographically secure random bytes. ERR: ", err)
		return "", errors.New("unable to generate random string")
	}

	return hex.EncodeToString(bytes), nil
}

func (r *CryptoRepository) HashSaltString(plaintext string) (string, error) {
	hash, err := argon2id.CreateHash(plaintext, argon2id.DefaultParams)
	if err != nil {
		r.logger.Error("Unable to hash password. ERR: ", err)
		return "", errors.New("unable to hash credential")
	}

	return hash, nil
}

func (r *CryptoRepository) IsPassHashMatching(plaintext string, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(plaintext, hash)
	if err != nil {
		r.logger.Error("Unable to compare password hash. ERR: ", err)
		return false, errors.New("unable to verify password")
	}

	return match, nil

}
