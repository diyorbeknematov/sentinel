package models

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	Id         uuid.UUID `json:"id"`
	UserId     string    `json:"user_id"`
	Type       string    `json:"type"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	IPAddress  string    `json:"ip_address"`
	RecordedAt time.Time `json:"recorded_at"`

	AgentId uuid.UUID `json:"agent_id"`
}

type CreateAppLog struct {
	AgentId   uuid.UUID `json:"agent_id" validate:"required"`
	UserId    string    `json:"user_id" validate:"required"`
	Type      string    `json:"type" validate:"required"`
	Level     string    `json:"level" validate:"required"`
	Message   string    `json:"message" validate:"required"`
	IPAddress string    `json:"ip_address" validate:"required"`
}

type FilterAppLog struct {
	AgentId   uuid.UUID `json:"agent_id"`
	UserId    string    `json:"user_id"`
	Type      string    `json:"type"`
	Level     string    `json:"level"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
