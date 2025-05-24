package entity

import (
	"time"

	"github.com/estella-studio/leon-backend/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	Email             string         `json:"email" gorm:"type:nvarchar(256);not null;unique"`
	Username          string         `json:"username" gorm:"type:nvarchar(64);not null;unique"`
	Password          string         `json:"password" gorm:"type:text;not null"`
	Name              string         `json:"name" gorm:"type:nvarchar(128)"`
	CreatedAt         time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	Data              []Data
	PasswordChange    []PasswordChange
	PasswordResetCode []PasswordResetCode
}

type PasswordChange struct {
	ID                uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID            uuid.UUID `json:"user_id" gorm:"type:char(36)"`
	Success           bool      `json:"success" gorm:"type:boolean"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	PasswordResetCode PasswordResetCode
}

type PasswordResetCode struct {
	ID               uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	PasswordChangeID uuid.UUID      `json:"change_id" gorm:"type:char(36)"`
	UserID           uuid.UUID      `json:"user_id" gorm:"type:char(36)"`
	Code             uint           `json:"code" gorm:"type:varchar(8)"`
	CreatedAt        time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (u *User) ParseToDTOResponseRegister() dto.ResponseRegister {
	return dto.ResponseRegister{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) ParseToDTOResponseLogin() dto.ResponseLogin {
	return dto.ResponseLogin{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) ParseToDTOResponseGetUserInfo() dto.ResponseGetUserInfo {
	return dto.ResponseGetUserInfo{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) ParseToDTOResponseUpdateUserInfo() dto.ResponseUpdateUserInfo {
	return dto.ResponseUpdateUserInfo{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
