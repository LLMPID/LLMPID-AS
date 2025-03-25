package service

import (
	"github.com/sirupsen/logrus"
	"llm-promp-inj.api/internal/repository"
)

type ExternalSystemService struct {
	CryptoRepo *repository.CryptoRepository
	UserRepo   *repository.UserRepository
	logger     *logrus.Logger
}

func NewExternalSystemService(cryptoRepo *repository.CryptoRepository, logger *logrus.Logger) *ExternalSystemService {
	return &ExternalSystemService{CryptoRepo: cryptoRepo, logger: logger}
}

func (s *ExternalSystemService) Register(systemName string) (string, error) {
	systemAccessKey, err := s.CryptoRepo.GenrateRandomString(32)
	if err != nil {
		return "", err
	}

	// The system access key is treated as a password.
	// Therefore, we hash+salt it.
	systemKeyHash, err := s.CryptoRepo.HashSaltString(systemAccessKey)
	if err != nil {
		return "", err
	}

	// The external system is treated as user with a role "ext_sys" internally.
	// Therefore, it is inserted in the users table.
	_, err = s.UserRepo.InsertUser(systemName, systemKeyHash, "ext_sys")
	if err != nil {
		return "", err
	}

	return systemAccessKey, nil
}
