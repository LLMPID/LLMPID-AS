package models

import "time"

type RevokedToken struct {
	ID        uint      `json:"id"`
	Token     string    `json:"token"`
	ExpiresAt int64     `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
