package dto

type AuthExtSystemRequest struct {
	SystemName string `json:"system_name" binding:"required" validate:"required,min=4,max=32"`
	AccessKey  string `json:"access_key" binding:"required" validate:"required,min=31"`
}
