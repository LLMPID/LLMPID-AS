package models

import (
	"time"
)

type ClassificationLog struct {
	ID          uint      `json:"id"`
	RequestText string    `json:"request_text"`
	Result      string    `json:"result"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
