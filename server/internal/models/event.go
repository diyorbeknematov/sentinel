package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type EventType string

const (
	EventMetric   EventType = "metrics"
	EventAppLog   EventType = "app_log"
	EventNginxLog EventType = "nginx_log"
)

type Event struct {
	Type    EventType       `json:"type"`
	AgentID uuid.UUID       `json:"agent_id"`
	Payload json.RawMessage `json:"payload"` // ← raw bytes saqlanadi
}
