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

type RawLog map[string]any
var TimeFormats = []string{
	time.RFC3339,
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
}

type MetricPayload struct {
	CPU     float64   `json:"cpu"`
	RAM     float64   `json:"ram"`
	Disk    float64   `json:"disk"`
	LogTime time.Time `json:"log_time"`
}

type AppLogPayload struct {
	ServiceName string         `json:"service_name"`
	Level       string         `json:"level"`
	Event       string         `json:"event"`
	Message     string         `json:"message"`
	LogTime     time.Time      `json:"log_time"`
	Metadata    map[string]any `json:"metadata,omitempty"`
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
