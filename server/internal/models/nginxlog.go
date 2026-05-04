package models

import (
	"time"

	"github.com/google/uuid"
)

type NginxLog struct {
	Id         uuid.UUID
	AgentId    uuid.UUID
	IPAddress  string
	Method     string
	Path       string
	Status     int
	Bytes      int
	UserAgent  string
	LogTime    time.Time
	RecordedAt time.Time
}

type CreateNginxLog struct {
	Id        uuid.UUID `db:"id"`
	AgentId   uuid.UUID `db:"agent_id"`
	IPAddress string    `db:"ip_address"`
	Method    string    `db:"method"`
	Path      string    `db:"path"`
	Status    int       `db:"status"`
	Bytes     int       `db:"bytes"`
	UserAgent string    `db:"user_agent"`
	LogTime   time.Time `db:"log_time"`
}

type NginxLogResponse struct {
	Id         uuid.UUID `json:"id"`
	AgentName  string    `json:"agent_name"`
	IPAddress  string    `json:"ip_address"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	Status     int       `json:"status"`
	Bytes      int       `json:"bytes"`
	UserAgent  string    `json:"user_agent"`
	LogTime    time.Time `json:"log_time"`
	RecordedAt time.Time `json:"recorded_at"`
}

type FilterNginxLog struct {
	AgentId string    `form:"agent_id"`
	Method  string    `form:"method"`
	Status  int       `form:"status"`
	From    time.Time `form:"from"`
	To      time.Time `form:"to"`
	Limit   int       `form:"limit"`
	Offset  int       `form:"-"`
}

type FilterNginxLogDB struct {
	AgentId uuid.UUID `db:"agent_id"`
	Method  string    `db:"method"`
	Status  int       `db:"status"`
	From    time.Time `db:"from"`
	To      time.Time `db:"to"`
	Limit   int       `db:"limit"`
	Offset  int       `db:"offset"`
}
