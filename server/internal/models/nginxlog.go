package models

import (
	"time"

	"github.com/google/uuid"
)

type NginxLog struct {
	Id         uuid.UUID `json:"id"`
	IPAddress  string    `json:"ip_address"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	Status     int       `json:"status"`
	Bytes      int       `json:"bytes"`
	UserAgent  string    `json:"user_agent"`
	LogTime    time.Time `json:"log_time"`
	RecordedAt time.Time `json:"recorded_at"`

	AgentId uuid.UUID `json:"agent_id"`
}

type CreateNginxLog struct {
	AgentId      uuid.UUID `json:"agent_id"`
	IPAddress    string    `json:"ip_address"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	Status       int       `json:"status"`
	Bytes      int       `json:"bytes"`
	UserAgent  string    `json:"user_agent"`
	LogTime    time.Time `json:"log_time"`
}

type FilterNginxLog struct {
	AgentId uuid.UUID `form:"agent_id"`
	Method  string    `form:"method"`
	Status  int       `form:"status"`
	From    time.Time `form:"from"`
	To      time.Time `form:"to"`

	Limit  int `form:"limit"`
	Offset int `form:"-"`
}
