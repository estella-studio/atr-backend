package dto

import (
	"time"

	"github.com/google/uuid"
)

type Register struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email" validate:"required,email"`
	Username string    `json:"username" validate:"required,min=4,max=20"`
	Password string    `json:"password" validate:"required,min=4"`
	Name     string    `json:"name" validate:"omitempty,min=3,max=28"`
}

type Login struct {
	Username string `json:"username" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=4"`
}

type UpdateUserInfo struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty,min=4,max=20"`
	Password string `json:"password" validate:"omitempty,min=4"`
	Name     string `json:"name" validate:"omitempty,min=3,max=28"`
}

type ResetPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordWithID struct {
	ID uuid.UUID `json:"id" validate:"required,min=36,max=36"`
}

type ResetPasswordWithCode struct {
	Email            string    `json:"email" validate:"required,email"`
	Code             uint      `json:"code" validate:"required"`
	Password         string    `json:"password" validate:"required,min=4"`
	PasswordChangeId uuid.UUID `json:"password_change_id"`
}

type ChangePassword struct {
	Password string `json:"password" validate:"required,min=4"`
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
