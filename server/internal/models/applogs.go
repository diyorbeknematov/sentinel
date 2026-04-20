package models

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	Id         uuid.UUID `json:"id"`
	UserId     string    `json:"user_id"`
	Event      string    `json:"event"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	LogTime    time.Time `json:"log_time"`
	RecordedAt time.Time `json:"recorded_at"`

	AgentId uuid.UUID `json:"agent_id"`
}

type CreateAppLog struct {
	AgentId uuid.UUID `json:"agent_id" validate:"required"`
	UserId  string    `json:"user_id" validate:"required"`
	Event   string    `json:"event" validate:"required"`
	Level   string    `json:"level" validate:"required"`
	Message string    `json:"message" validate:"required"`
	LogTime time.Time `json:"log_time" validate:"required"`
}

type FilterAppLog struct {
	AgentId uuid.UUID `json:"agent_id"`
	UserId  string    `json:"user_id"`
	Event   string    `json:"event"`
	Level   string    `json:"level"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
