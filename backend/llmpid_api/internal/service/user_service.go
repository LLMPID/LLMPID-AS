package service

import (
	"llm-promp-inj.api/internal/models"
	"llm-promp-inj.api/internal/repository"
)

type UserService struct {
	UserRepo   *repository.UserRepository
	CryptoRepo *repository.CryptoRepository
}

func NewUserService(userRepo *repository.UserRepository, cryptoRepo *repository.CryptoRepository) *UserService {
	return &UserService{UserRepo: userRepo, CryptoRepo: cryptoRepo}
}

func (s *UserService) Create(username string, password string, role string) (models.User, error) {
	var user models.User

	passwordHash, err := s.CryptoRepo.HashSaltString(password)
	if err != nil {
		return user, err
	}

	user, err = s.UserRepo.InsertUser(username, passwordHash, role)
	if err != nil {
		return user, err
	}

	return user, nil
}
