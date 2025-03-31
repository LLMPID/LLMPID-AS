package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"llm-promp-inj.api/internal/models"
)

type TokenRepository struct {
	serverSecretKey string
	DB              *gorm.DB
	logger          *logrus.Logger
}

func NewTokenRepository(serverSecretKey string, db *gorm.DB, logger *logrus.Logger) *TokenRepository {
	return &TokenRepository{serverSecretKey: serverSecretKey, DB: db, logger: logger}
}

func (r *TokenRepository) ExtractAndValidateJWT(tokenString string) (*models.AccessTokenClaims, error) {
	// Parse and verify the JWT token
	token, err := jwt.ParseWithClaims(tokenString, &models.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC and using HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(r.serverSecretKey), nil
	})

	// Handle parsing or verification errors
	if err != nil {
		r.logger.Error("Invalid signiture. ERR: ", err)
		return nil, err
	}

	// Extract token claims
	claims, ok := token.Claims.(*models.AccessTokenClaims)
	if !ok || !token.Valid {
		r.logger.Info("Token claims are invalid. ", claims)
		return nil, errors.New("invalid token")
	}

	// Perform additional validation for the token parameters
	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		r.logger.Info("Token has expired. ")
		return nil, errors.New("token has expired")
	}
	if claims.Issuer != "llmpid-api-service" {
		r.logger.Info("Token issuer is invalid. Issuer: ", claims.Issuer)
		return nil, errors.New("invalid token issuer")
	}

	return claims, nil
}

func (r *TokenRepository) GenerateJWT(username string, userID uint, expiration int64, role string, sessionSlug string) (string, models.AccessTokenClaims, error) {
	now := time.Now()

	sub := r.GenerateSubject(username, userID)
	accessClaims := models.AccessTokenClaims{
		Sub: sub,
		Data: map[string]string{
			"username": username,
			"role":     role,
		},
		SessionID: sessionSlug,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * time.Duration(expiration))),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "llmpid-api-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	tokenString, err := token.SignedString([]byte(r.serverSecretKey))
	if err != nil {
		r.logger.Error("Unable to sign user token. ERR: ", err)
		return "", accessClaims, errors.New("unable to create token")
	}

	return tokenString, accessClaims, nil
}

func (r *TokenRepository) GenerateSubject(username string, userID uint) string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s-%d", username, userID)))
	return hex.EncodeToString(hasher.Sum(nil)) // Encrypted-looking subject
}
