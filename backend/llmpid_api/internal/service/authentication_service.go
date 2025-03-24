package service

import (
	"errors"

	"llm-promp-inj.api/internal/models"
	"llm-promp-inj.api/internal/repository"
)

type AuthenticationService struct {
	UserRepo   *repository.UserRepository
	TokenRepo  *repository.TokenRepository
	CryptoRepo *repository.CryptoRepository
}

func NewAuthenticationService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository, cryptoRepo *repository.CryptoRepository) *AuthenticationService {
	return &AuthenticationService{
		UserRepo:   userRepo,
		TokenRepo:  tokenRepo,
		CryptoRepo: cryptoRepo,
	}
}

func (s *AuthenticationService) AuthenticateUser(username string, password string, sessionExpiration int) (string, error) {
	user, err := s.UserRepo.SelectUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	isValidPass, err := s.CryptoRepo.IsPassHashMatching(password, user.PasswordHash)
	if err != nil {
		return "", err
	}
	if !isValidPass {
		return "", nil
	}

	accessToken, _, err := s.TokenRepo.GenerateJWT(user.Username, user.ID, sessionExpiration, user.Role)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *AuthenticationService) RevokeUserSession(tokenString string) error {
	tokenClaims, err := s.TokenRepo.ExtractAndValidateJWT(tokenString)
	if err != nil {
		return err
	}

	revokedToken := models.RevokedToken{
		Token:     tokenString,
		ExpiresAt: tokenClaims.ExpiresAt.Unix(),
	}

	return s.TokenRepo.RevokeToken(revokedToken)
}
