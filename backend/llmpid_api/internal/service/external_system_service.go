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

func NewExternalSystemService(cryptoRepo *repository.CryptoRepository, userRepo *repository.UserRepository, logger *logrus.Logger) *ExternalSystemService {
	return &ExternalSystemService{CryptoRepo: cryptoRepo, UserRepo: userRepo, logger: logger}
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

func (s *ExternalSystemService) List() ([]string, error) {
	var services []string

	serviceUsers, err := s.UserRepo.SelectUserByRole("ext_sys")
	if err != nil {
		return []string{}, err
	}

	for _, serviceUser := range serviceUsers {
		services = append(services, serviceUser.Username)
	}

	return services, nil
}

func (s *ExternalSystemService) DeleteBySysName(username string) error {
	return s.UserRepo.DeleteByUsername(username)
}
