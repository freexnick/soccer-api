package response

import "github.com/google/uuid"

type UserResponse struct {
	UserID uuid.UUID `json:"id"`
	Email  string    `json:"email"`
	Token  string    `json:"token"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type RegisterResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
