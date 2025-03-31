package models

import "time"

type Session struct {
	ID        uint      `json:"id"`
	Sub       string    `json:"sub"`
	SessionID string    `json:"session_id"`
	ExpiresAt int64     `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
