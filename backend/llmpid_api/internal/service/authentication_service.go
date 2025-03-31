package service

import (
	"errors"

	"llm-promp-inj.api/internal/repository"
)

type AuthenticationService struct {
	UserRepo    *repository.UserRepository
	TokenRepo   *repository.TokenRepository
	CryptoRepo  *repository.CryptoRepository
	SessionRepo *repository.SessionRepository
}

func NewAuthenticationService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository, cryptoRepo *repository.CryptoRepository, sessionRepo *repository.SessionRepository) *AuthenticationService {
	return &AuthenticationService{
		UserRepo:    userRepo,
		TokenRepo:   tokenRepo,
		CryptoRepo:  cryptoRepo,
		SessionRepo: sessionRepo,
	}
}

func (s *AuthenticationService) Authenticate(username string, password string, sessionLength int64) (string, error) {
	// Retrieve the user from the DB in order to get the user's password hash.
	user, err := s.UserRepo.SelectUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Verifies whether the hash and the provided password match in order to authenticate the user.
	isValidPass, err := s.CryptoRepo.IsPassHashMatching(password, user.PasswordHash)
	if err != nil {
		return "", err
	}
	if !isValidPass {
		return "", nil
	}

	// Generate user session ID (SID) so sessions can be tracked and revoked.
	sessioID, _ := s.CryptoRepo.GenrateRandomString(32)

	// Generate a new access token for the user.
	accessToken, claims, err := s.TokenRepo.GenerateJWT(
		user.Username,
		user.ID,
		sessionLength,
		user.Role,
		sessioID)
	if err != nil {
		return "", err
	}

	// Create session for the user.
	// It can be revoked at any time - all tokens containing the session ID (sessionSlug) will be invalidated.
	err = s.SessionRepo.CreateSession(sessioID, claims.Sub, claims.ExpiresAt.Unix())
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *AuthenticationService) ChangePassword(username string, oldPassword string, newPassword string, tokenString string) (string, error) {
	// Retrieve the user from the DB in order to get the user's password hash.
	user, err := s.UserRepo.SelectUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Check if the old password provided by the user matches the one in the DB.
	// Used to verify that the user's identity is real and not a stolen access token.
	isValidPass, err := s.CryptoRepo.IsPassHashMatching(oldPassword, user.PasswordHash)
	if err != nil {
		return "", err
	}
	if !isValidPass {
		return "", nil
	}

	// Create a new argon2 hash (and salt) for the new password.
	newHash, err := s.CryptoRepo.HashSaltString(newPassword)
	if err != nil {
		return "", err
	}

	// Update the user's database entry with the new password hash.
	err = s.UserRepo.UpdatePasswordHashByUserID(user.ID, newHash)
	if err != nil {
		return "", err
	}

	// Revoke old user sessions.
	err = s.RevokeAllSessions(tokenString)
	if err != nil {
		return "", err
	}

	// Generate a new user session.
	newAccessToken, err := s.Authenticate(username, newPassword, 60)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func (s *AuthenticationService) IsValidSession(sessioID string, sub string) bool {
	return s.SessionRepo.IsValidSession(sessioID, sub)
}

func (s *AuthenticationService) RevokeSession(tokenString string) error {
	tokenClaims, err := s.TokenRepo.ExtractAndValidateJWT(tokenString)
	if err != nil {
		return err
	}

	s.SessionRepo.DeleteSessionBySID(tokenClaims.SessionID)
	return nil
}

func (s *AuthenticationService) RevokeAllSessions(tokenString string) error {
	tokenClaims, err := s.TokenRepo.ExtractAndValidateJWT(tokenString)
	if err != nil {
		return err
	}

	s.SessionRepo.DeleteSessionBySub(tokenClaims.Sub)
	return nil
}
