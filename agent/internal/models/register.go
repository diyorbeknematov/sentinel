package models

import "github.com/google/uuid"

type RegisterRequest struct {
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}

type RegisterResponse struct {
	AgentID      uuid.UUID `json:"id"`
	KafkaBrokers []string  `json:"kafka_brokers"`
	KafkaTopic   string    `json:"kafka_topic"`
}

type HeartbeatRequest struct {
	AgentID string `json:"agent_id"`
}
