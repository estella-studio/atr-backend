package mysql

import (
	"github.com/estella-studio/leon-backend/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		entity.User{},
		entity.UserDetail{},
		entity.Friend{},
		entity.FriendRequest{},
		entity.Verification{},
		entity.PasswordChange{},
		entity.PasswordResetCode{},
		entity.UserReporting{},
		entity.Data{},
	)

	return err
}
