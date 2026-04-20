package models

import (
	"time"

	"github.com/google/uuid"
)

type Metric struct {
	Id         uuid.UUID `json:"id"`
	CPU        float64   `json:"cpu"`
	RAM        float64   `json:"ram"`
	Disk       float64   `json:"disk"`
	LogTime    time.Time `json:"log_time"`
	RecordedAt time.Time `json:"recorded_at"`

	AgentId uuid.UUID `json:"agent_id"`
}

type CreateMetric struct {
	AgentId uuid.UUID `json:"agent_id" validate:"required"`
	CPU     float64   `json:"cpu" validate:"required"`
	RAM     float64   `json:"ram" validate:"required"`
	Disk    float64   `json:"disk" validate:"required"`
	LogTime time.Time `json:"log_time" validate:"required"`
}

type FilterMetrics struct {
	AgentId uuid.UUID `json:"agent_id"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
