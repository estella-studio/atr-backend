package dto

import (
	"time"

	"github.com/google/uuid"
)

type Add struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Type   bool      `json:"type" validate:"required,boolean"`
	Data   string    `json:"data"`
}

type Retrieve struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

type List struct {
	UserID uuid.UUID `json:"user_id"`
}

type ResponseAdd struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      bool      `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type ResponseRetrieve struct {
	Type bool   `json:"type"`
	Data string `json:"data"`
}

type ResponseList struct {
	ID        uuid.UUID `json:"id"`
	Type      bool      `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
