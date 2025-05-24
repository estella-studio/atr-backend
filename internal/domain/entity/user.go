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
	UserDetail        UserDetail
}

type UserDetail struct {
	ID           uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:char(36)"`
	ProfileIndex uint      `json:"profile_index" gorm:"type:tinyint unsigned"`
}

type Verification struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Email     string    `json:"email" gorm:"type:nvarchar(256);not null;unique"`
	Code      uint      `json:"code" gorm:"type:varchar(8)"`
	Success   bool      `json:"success" gorm:"type:boolean"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
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
	var responseRegister dto.ResponseRegister

	responseRegister.ID = u.ID
	responseRegister.Email = u.Email
	responseRegister.Username = u.Username
	responseRegister.Name = u.Name
	responseRegister.CreatedAt = u.CreatedAt
	responseRegister.UpdatedAt = u.UpdatedAt
	responseRegister.UserDetail.ProfileIndex = u.UserDetail.ProfileIndex

	return responseRegister
}

func (u *User) ParseToDTOResponseLogin() dto.ResponseLogin {
	var responseLogin dto.ResponseLogin

	responseLogin.ID = u.ID
	responseLogin.Email = u.Email
	responseLogin.Username = u.Username
	responseLogin.Name = u.Name
	responseLogin.CreatedAt = u.CreatedAt
	responseLogin.UpdatedAt = u.UpdatedAt
	responseLogin.UserDetail.ProfileIndex = u.UserDetail.ProfileIndex

	return responseLogin
}

func (u *User) ParseToDTOResponseGetUserInfo() dto.ResponseGetUserInfo {
	var responseGetUserInfo dto.ResponseGetUserInfo

	responseGetUserInfo.ID = u.ID
	responseGetUserInfo.Email = u.Email
	responseGetUserInfo.Username = u.Username
	responseGetUserInfo.Name = u.Name
	responseGetUserInfo.CreatedAt = u.CreatedAt
	responseGetUserInfo.UpdatedAt = u.UpdatedAt
	responseGetUserInfo.UserDetail.ProfileIndex = u.UserDetail.ProfileIndex

	return responseGetUserInfo
}

func (u *User) ParseToDTOResponseUpdateUserInfo() dto.ResponseUpdateUserInfo {
	var responseUdpateUserInfo dto.ResponseUpdateUserInfo

	responseUdpateUserInfo.ID = u.ID
	responseUdpateUserInfo.Email = u.Email
	responseUdpateUserInfo.Username = u.Username
	responseUdpateUserInfo.Name = u.Name
	responseUdpateUserInfo.CreatedAt = u.CreatedAt
	responseUdpateUserInfo.UpdatedAt = u.UpdatedAt
	responseUdpateUserInfo.UserDetail.ProfileIndex = u.UserDetail.ProfileIndex

	return responseUdpateUserInfo
}
