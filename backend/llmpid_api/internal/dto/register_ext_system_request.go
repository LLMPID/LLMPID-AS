package dto

type RegisterExtSystemRequest struct {
	SystemName string `json:"system_name" binding:"required" validate:"required,min=4,max=32"`
}
