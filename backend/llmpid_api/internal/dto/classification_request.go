package dto

type ClassificationRequest struct {
	Text string `json:"text" binding:"required"`
}
