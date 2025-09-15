package repository

import (
	"github.com/estella-studio/leon-backend/internal/domain/dto"
	"github.com/estella-studio/leon-backend/internal/domain/entity"
	"gorm.io/gorm"
)

type DataMySQLItf interface {
	Add(data *entity.Data) error
	Retrieve(data *entity.Data, userParam dto.Retrieve) error
	List(data *[]entity.Data, userParam dto.List) error
	ListPaged(data *[]entity.Data, userParam dto.List, offset int, limit int) error
	ListPublic(data *[]entity.Data, userParam dto.List) error
	ListPublicPaged(data *[]entity.Data, userParam dto.List, offset int, limit int) error
}

type DataMySQL struct {
	db *gorm.DB
}

func NewDataMySQL(db *gorm.DB) DataMySQLItf {
	return &DataMySQL{
		db: db,
	}
}

func (r *DataMySQL) Add(data *entity.Data) error {
	return r.db.Debug().
		Create(data).
		Error
}

func (r *DataMySQL) Retrieve(data *entity.Data, userParam dto.Retrieve) error {
	return r.db.Debug().
		Select("type, data").
		First(data, userParam).
		Error
}

func (r *DataMySQL) List(data *[]entity.Data, userParam dto.List) error {
	return r.db.Debug().
		Select("id, type, created_at").
		Find(data, userParam).
		Error
}

func (r *DataMySQL) ListPaged(data *[]entity.Data, userParam dto.List, offset int, limit int) error {
	return r.db.Debug().
		Select("id, type, created_at").
		Limit(limit).
		Offset(offset).
		Find(data, userParam).
		Error
}

func (r *DataMySQL) ListPublic(data *[]entity.Data, userParam dto.List) error {
	return r.db.Debug().
		Select("id, type, created_at").
		Where("type = ? ", true).
		Find(data, userParam).
		Error
}

func (r *DataMySQL) ListPublicPaged(data *[]entity.Data, userParam dto.List, offset int, limit int) error {
	return r.db.Debug().
		Select("id, type, created_at").
		Where("type = ? ", true).
		Limit(limit).
		Offset(offset).
		Find(data, userParam).
		Error
}
