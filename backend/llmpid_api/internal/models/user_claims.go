package models

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaims struct {
	Sub       string            `json:"sub"`
	Data      map[string]string `json:"data"`
	SessionID string            `json:"session_id"`
	jwt.RegisteredClaims
}
