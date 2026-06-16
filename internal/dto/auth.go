package dto

import types "godance/internal/type"

type RegistrationRequest struct {
	Email    string         `json:"email" binding:"required"`
	Phone    string         `json:"phone" binding:"required"`
	Password string         `json:"password" binding:"required"`
	Role     types.UserRole `json:"role" binding:"required"`
	Fullname string         `json:"full_name" binding:"required"`
	Bio      *string        `json:"bio" binding:"omitempty"`
}

type AuthResponse struct {
	Access_token  string    `json:"access_token"`
	Refresh_token string    `json:"refresh_token"`
	User          UserShort `json:"user"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
