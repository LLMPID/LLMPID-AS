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
	// Check if the token has been blaclisted by a logout event
	if err := r.DB.Where("token = ?", tokenString).First(&models.RevokedToken{}).Error; err == nil {
		return nil, nil
	}

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

func (r *TokenRepository) GenerateJWT(username string, userID uint, expiration int, role string) (string, models.AccessTokenClaims, error) {
	now := time.Now()

	accessClaims := models.AccessTokenClaims{
		Sub: r.GenerateSubject(username, userID),
		Data: map[string]string{
			"role": role,
		},
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

func (r *TokenRepository) RevokeToken(token models.RevokedToken) error {
	err := r.DB.Create(&token).Error
	if err != nil {
		r.logger.Error("Unable to insert revoked token. ERR: ", err)
		return errors.New("unable to revoke token")
	}

	return nil
}

func (r *TokenRepository) GenerateSubject(username string, userID uint) string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s-%d", username, userID)))
	return hex.EncodeToString(hasher.Sum(nil)) // Encrypted-looking subject
}
