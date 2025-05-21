package mysql

import (
	"github.com/estella-studio/leon-backend/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		entity.User{},
		entity.PasswordChange{},
		entity.Data{},
	)

	return err
}
