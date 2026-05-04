package models

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	Id         uuid.UUID
	AgentId    uuid.UUID
	UserId     string
	Event      string
	Level      string
	Message    string
	LogTime    time.Time
	RecordedAt time.Time
}

type CreateAppLog struct {
	Id      uuid.UUID `db:"id"`
	AgentId uuid.UUID `db:"agent_id"`
	UserId  string    `db:"user_id"`
	Event   string    `db:"event"`
	Level   string    `db:"level"`
	Message string    `db:"message"`
	LogTime time.Time `db:"log_time"`
}

type AppLogResponse struct {
	Id         uuid.UUID `json:"id"`
	AgentName  string    `json:"agent_name"`
	UserId     string    `json:"user_id,omitempty"`
	Event      string    `json:"event"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	LogTime    time.Time `json:"log_time"`
	RecordedAt time.Time `json:"recorded_at"`
}

type FilterAppLog struct {
	AgentId string    `form:"agent_id"`
	UserId  string    `form:"user_id"`
	Event   string    `form:"event"`
	Level   string    `form:"level"`
	From    time.Time `form:"from"`
	To      time.Time `form:"to"`
	Limit   int       `form:"limit"`
	Offset  int       `form:"-"`
}

type FilterAppLogDB struct {
	AgentId uuid.UUID `db:"agent_id"`
	UserId  string    `db:"user_id"`
	Event   string    `db:"event"`
	Level   string    `db:"level"`
	From    time.Time `db:"from"`
	To      time.Time `db:"to"`
	Limit   int       `db:"limit"`
	Offset  int       `db:"offset"`
}
