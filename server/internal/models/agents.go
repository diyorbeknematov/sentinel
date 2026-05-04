package models

import (
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	Id        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status,omitempty"`
	IPAddress string    `json:"ip_address"`
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAgent struct {
	Id        string `json:"-"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}

type CreateAgentDB struct {
	Id        uuid.UUID `db:"id"`
	AccountID uuid.UUID `db:"account_id"`
	Name      string    `db:"name"`
	IPAddress string    `db:"ip_address"`
	LastSeen  time.Time `db:"last_seen"`
}

type CreateAgentResponse struct {
	ID           uuid.UUID `json:"id"`
	KafkaBrokers []string  `json:"kafka_brokers"`
	KafkaTopic   string    `json:"kafka_topic"`
}

type HeartbeatRequest struct {
	AgentID uuid.UUID `json:"agent_id"`
}

type UpdateLastSeen struct {
	Id       uuid.UUID `db:"id"`
	LastSeen time.Time `db:"last_seen"`
}

type FilterAgentDB struct {
	AccountID uuid.UUID `db:"account_id"`
	Name      string    `db:"name"`
	From      time.Time `db:"from"`
	To        time.Time `db:"to"`
	Limit     int       `db:"limit"`
	Offset    int       `db:"offset"`
}

type FilterAgent struct {
	AccountID string    `form:"account_id"`
	Name      string    `form:"name"`
	From      time.Time `form:"from"`
	To        time.Time `form:"to"`
	Limit     int       `form:"limit"`
	Offset    int       `form:"-"`
}
