package models

import (
	"time"

	"github.com/google/uuid"
)

type Metric struct {
	Id         uuid.UUID `json:"id"`
	CPU        float32   `json:"cpu"`
	RAM        float32   `json:"ram"`
	Disk       float32   `json:"disk"`
	RecordedAt time.Time `json:"recorded_at"`

	AgentId uuid.UUID `json:"agent_id"`
}

type CreateMetric struct {
	AgentId uuid.UUID `json:"agent_id" validate:"required"`
	CPU     float32   `json:"cpu" validate:"required"`
	RAM     float32   `json:"ram" validate:"required"`
	Disk    float32   `json:"disk" validate:"required"`
}

type FilterMetrics struct {
	AgentId uuid.UUID    `json:"agent_id"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`

	Limit  int
	Offest int
}
