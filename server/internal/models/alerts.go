package models

import (
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	Id        uuid.UUID
	Type      string
	Message   string
	Severity  string
	CreatedAt time.Time
	IsRead    bool
	AgentId   uuid.UUID
}

type CreateAlert struct {
	Id       uuid.UUID `db:"id"`
	AgentId  uuid.UUID `db:"agent_id"`
	Type     string    `db:"type"`
	Message  string    `db:"message"`
	Severity string    `db:"severity"`
}

type AlertResponse struct {
	Id        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	CreatedAt time.Time `json:"created_at"`
	IsRead    bool      `json:"is_read"`
	AgentName string    `json:"agent_name"`
}

type FilterAlert struct {
	AgentId  *string    `form:"agent_id"`
	Severity *string    `form:"severity"`
	IsRead   *bool      `form:"is_read"`
	From     *time.Time `form:"from" time_format:"2006-01-02"`
	To       *time.Time `form:"to" time_format:"2006-01-02"`
	Limit    int       `form:"limit"`
	Offset   int       `form:"-"`
}

type FilterAlertDB struct {
	AgentId  *uuid.UUID `db:"agent_id"`
	Severity *string    `db:"severity"`
	IsRead   *bool      `db:"is_read"`
	From     *time.Time `db:"from"`
	To       *time.Time `db:"to"`
	Limit    int        `db:"limit"`
	Offset   int        `db:"offset"`
}
