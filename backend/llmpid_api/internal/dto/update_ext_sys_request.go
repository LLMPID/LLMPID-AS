package dto

type UpdateExtSystemRequest struct {
	OldSystemName string `json:"old_system_name" binding:"required" validate:"required,min=4,max=32"`
	NewSystemName string `json:"new_system_name" binding:"required" validate:"required,min=4,max=32"`
}
