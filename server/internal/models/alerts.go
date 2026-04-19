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
	Id     uuid.UUID `json:"-"`
	IsRead bool      `json:"is_read" validate:"required"`
}

type FilterAlert struct {
	AgentId  *uuid.UUID `form:"agent_id"`
	Type     *string    `form:"type"`
	Severity *string    `form:"severity"`
	IsRead   *bool      `form:"is_read"`

	From *time.Time `form:"from" time_format:"2006-01-02"`
	To   *time.Time `form:"to" time_format:"2006-01-02"`

	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}
