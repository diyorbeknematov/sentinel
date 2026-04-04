package models

import (
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	Id        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	CreatedAt time.Time `json:"created_at"`
	IsRead    bool      `json:"-"`

	AgentId uuid.UUID `json:"agent_id"`
}

type CreateAlert struct {
	AgentId  uuid.UUID `json:"agent_id" validate:"required"`
	Type     string    `json:"type" validate:"required"`
	Message  string    `json:"message" validate:"required"`
	Severity string    `json:"severity" validate:"required"`
}

type MarkAlertRead struct {
	Id     uuid.UUID `json:"id" validate:"required"`
	IsRead bool      `json:"is_read" validate:"required"`
}
