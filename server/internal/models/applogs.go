package models

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	Id          uuid.UUID      `json:"id"`
	AgentId     uuid.UUID      `json:"agent_id"`
	ServiceName string         `json:"service_name"`
	Event       string         `json:"event"`
	Level       string         `json:"level"`
	Message     string         `json:"message"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	LogTime     time.Time      `json:"log_time"`
	RecordedAt  time.Time      `json:"recorded_at"`
}

type CreateAppLog struct {
	Id          uuid.UUID `db:"id"`
	AgentId     uuid.UUID `db:"agent_id"`
	ServiceName string    `db:"service_name"`
	Event       string    `db:"event"`
	Level       string    `db:"level"`
	Message     string    `db:"message"`
	Metadata    string    `db:"metadata,omitempty"`
	LogTime     time.Time `db:"log_time"`
}

type AppLogResponse struct {
	Id          uuid.UUID      `json:"id"`
	AgentName   string         `json:"agent_name"`
	ServiceName string         `json:"service_name"`
	Event       string         `json:"event"`
	Level       string         `json:"level"`
	Message     string         `json:"message"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	LogTime     time.Time      `json:"log_time"`
	RecordedAt  time.Time      `json:"recorded_at"`
}

type FilterAppLog struct {
	AgentId string    `form:"agent_id"`
	Level   string    `form:"level"`
	From    time.Time `form:"from"`
	To      time.Time `form:"to"`
	Limit   int       `form:"limit"`
	Offset  int       `form:"-"`
}

type FilterAppLogDB struct {
	AgentId uuid.UUID `db:"agent_id"`
	Level   string    `db:"level"`
	From    time.Time `db:"from"`
	To      time.Time `db:"to"`
	Limit   int       `db:"limit"`
	Offset  int       `db:"offset"`
}
