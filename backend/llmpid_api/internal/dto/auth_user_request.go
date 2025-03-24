package dto

type AuthUserRequest struct {
	Usernames string `json:"username" binding:"required" validate:"required,min=8,max=32"`
	Password  string `json:"password" binding:"required" validate:"required,min=8"`
}
