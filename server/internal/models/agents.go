package models

import (
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	IPAddress string    `json:"ip_address"`
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAgent struct {
	Name      string    `json:"name" validate:"required"`
	APIKey    string    `json:"api_key" validate:"required"`
	IPAddress string    `json:"ip_address" validate:"required"`
	LastSeen  time.Time `json:"last_seen" validate:"required"`
}

type UpdateLastSeen struct {
	Id       uuid.UUID `json:"id" validate:"required"`
	LastSeen time.Time `json:"last_seen" validate:"required"`
}

type FilterAgent struct {
	Name      string    `json:"name"`
	IPAddress string    `json:"ip_address,omitempty"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
