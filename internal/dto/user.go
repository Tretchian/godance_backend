package dto

import types "godance/internal/type"

type UserShort struct {
	ID       uint           `json:"id"`
	Email    string         `json:"email"`
	Role     types.UserRole `json:"role"`
	Fullname string         `json:"full_name"`
}
