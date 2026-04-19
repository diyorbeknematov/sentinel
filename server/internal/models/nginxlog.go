package models

import (
	"time"

	"github.com/google/uuid"
)

type NginxLog struct {
	Id           uuid.UUID `json:"id"`
	IPAddress    string    `json:"ip_address"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	Status       int       `json:"status"`
	ResponseTime int       `json:"response_time"`
	RecordedAt   time.Time `json:"recorded_at"`

	AgentId uuid.UUID `json:"agent_id"`
}

type CreateNginxLog struct {
	AgentId      uuid.UUID `json:"agent_id"`
	IPAddress    string    `json:"ip_address"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	Status       int       `json:"status"`
	ResponseTime int       `json:"response_time"`
}

type FilterNginxLog struct {
	AgentId      uuid.UUID `json:"agent_id"`
	Method       string    `json:"method"`
	Status       int       `json:"status"`
	From         time.Time `json:"from"`
	To           time.Time `json:"to"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
