package entity

import (
	"time"

	"github.com/estella-studio/atr-backend/internal/domain/dto"
	"github.com/google/uuid"
)

type Data struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36)"`
	Data      string    `json:"data" gorm:"type:varchar(256)"`
	Type      bool      `json:"type" gorm:"type:boolean"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
}

func (d *Data) ParseToDTOResponseAdd() dto.ResponseAdd {
	return dto.ResponseAdd{
		ID:        d.ID,
		UserID:    d.UserID,
		Type:      d.Type,
		CreatedAt: d.CreatedAt,
	}
}

func (d *Data) ParseToDTOResponseRetrieve() dto.ResponseRetrieve {
	return dto.ResponseRetrieve{
		Type: d.Type,
		Data: d.Data,
	}
}

func (d *Data) ParseToDTOResponseList() dto.ResponseList {
	return dto.ResponseList{
		ID:        d.ID,
		Type:      d.Type,
		CreatedAt: d.CreatedAt,
	}
}
