package models

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	Id         uuid.UUID `json:"id"`
	Type       string    `json:"type"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	IPAddress  string    `json:"ip_address"`
	RecordedAt time.Time `json:"recorded_at"`

	AgentId uuid.UUID `json:"agent_id"`
}

type CreateLog struct {
	AgentId   uuid.UUID `json:"agent_id" validate:"required"`
	Type      string    `json:"type" validate:"required"`
	Level     string    `json:"level" validate:"required"`
	Message   string    `json:"message" validate:"required"`
	IPAddress string    `json:"ip_address" validate:"required"`
}