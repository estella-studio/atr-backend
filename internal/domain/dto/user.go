package dto

import (
	"time"

	"github.com/google/uuid"
)

type Register struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email" validate:"omitempty,email"`
	Username string    `json:"username" validate:"required,min=3,max=64"`
	Password string    `json:"password" validate:"required,min=8,max=256"`
	Name     string    `json:"name" validate:"omitempty,min=3,max=128"`
}

type Login struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserInfo struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty,min=3"`
	Password string `json:"password" validate:"omitempty,min=8"`
	Name     string `json:"name" validate:"omitempty,min=3"`
}

type ResponseRegister struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseLogin struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseGetUserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseUpdateUserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
