package repository

import (
	"github.com/estella-studio/leon-backend/internal/domain/dto"
	"github.com/estella-studio/leon-backend/internal/domain/entity"
	"gorm.io/gorm"
)

type UserMySQLItf interface {
	Register(user *entity.User) error
	Login(user *entity.User) error
	GetUserInfo(user *entity.User) error
	UpdateUserInfo(user *entity.User) error
	GetUsername(user *entity.User, userParam dto.Login) error
}

type UserMySQL struct {
	db *gorm.DB
}

func NewUserMySQL(db *gorm.DB) UserMySQLItf {
	return &UserMySQL{
		db,
	}
}

func (r *UserMySQL) Register(user *entity.User) error {
	return r.db.Debug().Create(user).Error
}

func (r *UserMySQL) Login(user *entity.User) error {
	return r.db.Debug().First(user).Error
}

func (r *UserMySQL) GetUserInfo(user *entity.User) error {
	return r.db.Debug().First(user).Error
}

func (r *UserMySQL) UpdateUserInfo(user *entity.User) error {
	return r.db.Debug().Updates(user).Error
}

func (r *UserMySQL) GetUsername(user *entity.User, userParam dto.Login) error {
	return r.db.Debug().First(&user, userParam).Error
}
