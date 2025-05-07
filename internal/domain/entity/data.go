package entity

import (
	"time"

	"github.com/estella-studio/leon-backend/internal/domain/dto"
	"github.com/google/uuid"
)

type Data struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36)"`
	Data      []byte    `json:"data" gorm:"type:blob"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
}

func (d *Data) ParseToDTOResponseAdd() dto.ResponseAdd {
	return dto.ResponseAdd{
		ID:        d.ID,
		UserID:    d.UserID,
		CreatedAt: d.CreatedAt,
	}
}

func (d *Data) ParseToDTOResponseRetrieve() dto.ResponseRetrieve {
	return dto.ResponseRetrieve{
		ID:        d.ID,
		UserID:    d.UserID,
		Data:      d.Data,
		CreatedAt: d.CreatedAt,
	}
}
