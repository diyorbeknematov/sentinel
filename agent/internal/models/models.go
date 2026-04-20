package models

import (
	"time"

	"github.com/google/uuid"
)

type AppLogPayload struct {
	AgentID uuid.UUID `json:"agent_id"` // UUID
	UserID  string    `json:"user_id"`
	Event   string    `json:"event"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
	LogTime time.Time `json:"timestamp"`
}

type NginxLogPayload struct {
	IP        string    `json:"ip"`
	AgentID   uuid.UUID `json:"agent_id"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Status    int       `json:"status"`
	Bytes     int       `json:"bytes"`
	UserAgent string    `json:"user_agent"`
	LogTime   time.Time `json:"time"`
}

type MetricPayload struct {
	AgentID uuid.UUID `json:"agent_id"`
	CPU     float64   `json:"cpu"`
	RAM     float64   `json:"ram"`
	Disk    float64   `json:"disk"`
	LogTime time.Time `json:"log_time"`
}
