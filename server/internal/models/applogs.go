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
	AgentId uuid.UUID `form:"agent_id"`
	UserId  string    `form:"user_id"`
	Event   string    `form:"event"`
	Level   string    `form:"level"`
	From    time.Time `form:"from"`
	To      time.Time `form:"to"`

	Limit  int `form:"limit"`
	Offset int `form:"-"`
}
