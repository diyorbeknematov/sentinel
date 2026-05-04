package models

import (
	"time"

	"github.com/google/uuid"
)

type Metric struct {
	Id         uuid.UUID
	AgentId    uuid.UUID
	CPU        float64
	RAM        float64
	Disk       float64
	LogTime    time.Time
	RecordedAt time.Time
}

type CreateMetric struct {
	Id      uuid.UUID `db:"id"`
	AgentId uuid.UUID `db:"agent_id"`
	CPU     float64   `db:"cpu"`
	RAM     float64   `db:"ram"`
	Disk    float64   `db:"disk"`
	LogTime time.Time `db:"log_time"`
}

type MetricResponse struct {
	Id         uuid.UUID `json:"id"`
	AgentName  string    `json:"agent_name"`
	CPU        float64   `json:"cpu"`
	RAM        float64   `json:"ram"`
	Disk       float64   `json:"disk"`
	LogTime    time.Time `json:"log_time"`
	RecordedAt time.Time `json:"recorded_at"`
}

type FilterMetrics struct {
	AgentId string    `form:"agent_id"`
	From    time.Time `form:"from"`
	To      time.Time `form:"to"`
	Limit   int       `form:"limit"`
	Offset  int       `form:"-"`
}

type FilterMetricsDB struct {
	AgentId uuid.UUID `db:"agent_id"`
	From    time.Time `db:"from"`
	To      time.Time `db:"to"`
	Limit   int       `db:"limit"`
	Offset  int       `db:"offset"`
}
