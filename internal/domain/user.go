package domain

import "time"

type User struct {
	ID        string    `json:"id" validate:"omitempty,uuid"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"-" validate:"required,min=8"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
