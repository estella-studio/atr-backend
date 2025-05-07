package repository

import (
	"github.com/estella-studio/leon-backend/internal/domain/entity"
	"gorm.io/gorm"
)

type DataMySQLItf interface {
	Add(data *entity.Data) error
	Retrieve(data *entity.Data) error
}

type DataMySQL struct {
	db *gorm.DB
}

func NewDataMySQL(db *gorm.DB) DataMySQLItf {
	return &DataMySQL{
		db,
	}
}

func (r *DataMySQL) Add(data *entity.Data) error {
	return r.db.Debug().Create(data).Error
}

func (r *DataMySQL) Retrieve(data *entity.Data) error {
	return r.db.Debug().First(data).Error
}
