package service

import (
	"llm-promp-inj.api/internal/models"
	"llm-promp-inj.api/internal/repository"
)

type TokenService struct {
	TokenRepo *repository.TokenRepository
}

func NewTokenService(tokenRepository *repository.TokenRepository) *TokenService {
	return &TokenService{TokenRepo: tokenRepository}
}

func (s *TokenService) ExtractAndValidateToken(tokenString string) (*models.AccessTokenClaims, error) {
	claims, err := s.TokenRepo.ExtractAndValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
