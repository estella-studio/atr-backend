package repository

import (
	"log"

	"github.com/estella-studio/leon-backend/internal/domain/dto"
	"github.com/estella-studio/leon-backend/internal/domain/entity"
	"gorm.io/gorm"
)

type UserMySQLItf interface {
	Register(user *entity.User) error
	Login(user *entity.User) error
	ChangePassword(user *entity.User) error
	GetEmail(user *entity.User, userParam dto.ResetPassword) error
	GetUserInfo(user *entity.User) error
	GetUserID(user *entity.User, userParam dto.ResetPassword) error
	GetUsername(user *entity.User, userParam dto.Login) error
	GetPasswordChangeID(passwordChange *entity.PasswordChange, userParam dto.ResetPassword) error
	GetPasswordResetCode(passwordResetcode *entity.PasswordResetCode, userParam dto.ResetPasswordWithCode) error
	GetPasswordChangeValidity(passwordChange *entity.PasswordChange) error
	GetPasswordResetCodeValidity(passwordChange *entity.PasswordResetCode) error
	GetPasswordChangeEntry(passwordChange *entity.PasswordChange, userParam dto.ResetPasswordWithID) error
	CreatePasswordChangeEntry(passwordChange *entity.PasswordChange) error
	CreatePasswordResetCode(passwordResetCode *entity.PasswordResetCode) error
	UpdatePasswordResetCode(passwordResetCode *entity.PasswordResetCode) error
	UpdatePasswordChangeEntry(passwordChange *entity.PasswordChange) error
	UpdateUserInfo(user *entity.User) error
	SoftDelete(user *entity.User) error
}

type UserMySQL struct {
	db *gorm.DB
}

func NewUserMySQL(db *gorm.DB) UserMySQLItf {
	return &UserMySQL{
		db: db,
	}
}

func (r *UserMySQL) Register(user *entity.User) error {
	return r.db.Debug().
		Create(user).
		Error
}

func (r *UserMySQL) Login(user *entity.User) error {
	return r.db.Debug().
		First(user).
		Error
}

func (r *UserMySQL) ChangePassword(user *entity.User) error {
	return r.db.Debug().
		Model(&user).
		Update("password", user.Password).
		Error
}

func (r *UserMySQL) GetEmail(user *entity.User, userParam dto.ResetPassword) error {
	return r.db.Debug().
		Select("email").
		First(user, userParam).
		Error
}

func (r *UserMySQL) GetUserInfo(user *entity.User) error {
	return r.db.Debug().
		Select("id", "email", "username", "name", "created_at", "updated_at").
		First(user).
		Error
}

func (r *UserMySQL) GetUserID(user *entity.User, userParam dto.ResetPassword) error {
	return r.db.Debug().
		Select("id").
		First(user, userParam).
		Error
}

func (r *UserMySQL) GetUsername(user *entity.User, userParam dto.Login) error {
	return r.db.Debug().
		First(user, userParam).
		Error
}

func (r *UserMySQL) GetPasswordChangeID(passwordChange *entity.PasswordChange, userParam dto.ResetPassword) error {
	return r.db.Debug().
		Select("id").
		First(passwordChange, userParam).
		Error
}

func (r *UserMySQL) GetPasswordResetCode(passwordResetcode *entity.PasswordResetCode, userParam dto.ResetPasswordWithCode) error {
	return r.db.Debug().
		Select("password_change_id", "code").
		First(passwordResetcode, userParam).
		Error
}

func (r *UserMySQL) GetPasswordResetCodeValidity(passwordChange *entity.PasswordResetCode) error {
	return r.db.Debug().
		Order("created_at desc").
		Find(&passwordChange).
		Select("created_at").
		Where("user_id = ?", passwordChange.UserID).
		Error
}

func (r *UserMySQL) GetPasswordChangeValidity(passwordChange *entity.PasswordChange) error {
	return r.db.Debug().
		Select("id", "created_at", "success").
		First(passwordChange).
		Error
}

func (r *UserMySQL) GetPasswordChangeEntry(passwordChange *entity.PasswordChange, userParam dto.ResetPasswordWithID) error {
	return r.db.Debug().
		Select("user_id").
		First(passwordChange).
		Error
}

func (r *UserMySQL) CreatePasswordChangeEntry(passwordChange *entity.PasswordChange) error {
	return r.db.Debug().
		Create(passwordChange).
		Error
}

func (r *UserMySQL) CreatePasswordResetCode(passwordResetCode *entity.PasswordResetCode) error {
	return r.db.Debug().
		Create(passwordResetCode).
		Error
}

func (r *UserMySQL) UpdatePasswordResetCode(passwordResetCode *entity.PasswordResetCode) error {
	return r.db.Debug().
		Model(&passwordResetCode).
		Where("password_change_id = ?", passwordResetCode.PasswordChangeID).
		Delete(passwordResetCode).
		Error
}

func (r *UserMySQL) UpdatePasswordChangeEntry(passwordChange *entity.PasswordChange) error {
	log.Println(passwordChange)
	return r.db.Debug().
		Model(&passwordChange).
		Where("id = ?", passwordChange.ID).
		Update("success", passwordChange.Success).
		Error
}

func (r *UserMySQL) UpdateUserInfo(user *entity.User) error {
	return r.db.Debug().
		Updates(user).
		Error
}

func (r *UserMySQL) SoftDelete(user *entity.User) error {
	return r.db.Debug().
		Delete(user).
		Error
}
