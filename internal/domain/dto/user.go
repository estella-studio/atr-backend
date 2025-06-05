package dto

import (
	"time"

	"github.com/google/uuid"
)

type Register struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email" validate:"required,email"`
	Username     string    `json:"username" validate:"required,min=4,max=20"`
	Password     string    `json:"password" validate:"required,min=4"`
	Name         string    `json:"name" validate:"omitempty,min=3,max=29"`
	ProfileIndex uint      `json:"profile_index" validate:"omitempty"`
}

type Login struct {
	Username string `json:"username" validate:"required,min=4,max=20"`
	Password string `json:"password" validate:"required,min=4"`
}

type RenewToken struct {
	ID uuid.UUID `json:"id"`
}

type UserDetail struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	ProfileIndex uint      `json:"profile_index" validate:"omitempty"`
}

type AddFriend struct {
	UserID   uuid.UUID `json:"user_id"`
	FriendID uuid.UUID `json:"friend_id" validate:"required,uuid_rfc4122"`
}

type CheckUserID struct {
	ID uuid.UUID `json:"id" validate:"required,uuid_rfc4122"`
}

type CheckFriendRequestExist struct {
	UserID   uuid.UUID `json:"user_id"`
	FriendID uuid.UUID `json:"friend_id"`
}

type AcceptFriendRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	FriendID uuid.UUID `json:"friend_id"`
	Username string    `json:"username" validate:"required,min=4,max=20"`
}

type SendFriendRequest struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	FriendID uuid.UUID `json:"friend_id"`
	Username string    `json:"username" validate:"required,min=4,max=20"`
	Accepted bool      `json:"accepted"`
}

type GetFriendRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	FriendID uuid.UUID `json:"friend_id"`
}

type GetUserID struct {
	Username string `json:"username" validate:"required,min=4,max=20"`
}

type GetUserInfoPublic struct {
	Username string `json:"username" validate:"required,min=4,max=20"`
}

type UpdateUserInfo struct {
	Email        string `json:"email" validate:"omitempty,email"`
	Username     string `json:"username" validate:"omitempty,min=4,max=20"`
	Password     string `json:"password" validate:"omitempty,min=4"`
	Name         string `json:"name" validate:"omitempty,min=3,max=29"`
	ProfileIndex uint   `json:"profile_index" validate:"omitempty"`
}

type EmailVerification struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email" validate:"required,email"`
	Code  uint      `json:"code"`
}

type ValidateEmail struct {
	Email string `json:"email" validate:"required,email"`
	Code  uint   `json:"code" validate:"required,min=8"`
}

type CheckUsername struct {
	Username string `json:"username" validate:"required,min=4,max=20"`
}

type ResetPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordWithID struct {
	ID uuid.UUID `json:"id" validate:"required,min=36,max=36"`
}

type CheckPasswordResetCode struct {
	Email string `json:"email" validate:"required,email"`
	Code  uint   `json:"code" validate:"required,min=8"`
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

type ReportUser struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Username   string    `json:"username" validate:"required,min=4,max=20"`
	ReporterID uuid.UUID `json:"reporter_id"`
}

type ResponseRegister struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserDetail struct {
		ProfileIndex uint `json:"profile_index"`
	} `json:"user_detail"`
}

type ResponseLogin struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserDetail struct {
		ProfileIndex uint `json:"profile_index"`
	} `json:"user_detail"`
}

type ResponseGetUserInfo struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserDetail struct {
		ProfileIndex uint      `json:"profile_index"`
		LastActivity time.Time `json:"last_activity"`
	} `json:"user_detail"`
}

type ResponseGetUserInfoPublic struct {
	Username   string `json:"username"`
	Name       string `json:"name"`
	UserDetail struct {
		ProfileIndex uint      `json:"profile_index"`
		LastActivity time.Time `json:"last_activity"`
	} `json:"user_detail"`
}

type ResponseUpdateUserInfo struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserDetail struct {
		ProfileIndex uint `json:"profile_index"`
	} `json:"user_detail"`
}

type ResponseGetFriendRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	FriendID uuid.UUID `json:"friend_id"`
}

type ResponseFriendList struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}
