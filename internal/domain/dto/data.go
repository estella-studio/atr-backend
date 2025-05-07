package dto

import (
	"time"

	"github.com/google/uuid"
)

type Add struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Data   []byte    `json:"data"`
}

type ResponseAdd struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ResponseRetrieve struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Data      []byte    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}
