package models

import "time"

type EventType string

const (
	EventMetric   EventType = "metrics"
	EventAppLog   EventType = "app_log"
	EventNginxLog EventType = "nginx_log"
)

type Event struct {
	Type    EventType `json:"type"`
	AgentID string    `json:"agent_id"`
	Payload any       `json:"payload"`
}

type MetricPayload struct {
	CPU     float64   `json:"cpu"`
	RAM     float64   `json:"ram"`
	Disk    float64   `json:"disk"`
	LogTime time.Time `json:"log_time"`
}

type AppLogPayload struct {
	UserId  string    `json:"user_id"`
	Level   string    `json:"level"`
	Event   string    `json:"event"`
	Message string    `json:"message"`
	LogTime time.Time `json:"log_time"`
}

type NginxLogPayload struct {
	IPAddress string    `json:"ip_address"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Status    int       `json:"status"`
	Bytes     int       `json:"bytes"`
	UserAgent string    `json:"user_agent"`
	LogTime   time.Time `json:"log_time"`
}
