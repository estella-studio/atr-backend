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

type Retrieve struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

type ResponseAdd struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ResponseRetrieve struct {
	Data []byte `json:"data"`
}

type ResponseList struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
