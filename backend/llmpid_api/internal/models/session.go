package models

import "time"

type Session struct {
	ID        uint      `json:"id"`
	Sub       string    `json:"sub"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
