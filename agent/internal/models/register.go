package models

import "github.com/google/uuid"

type RegisterRequest struct {
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}

type RegisterResponse struct {
	AgentId uuid.UUID `json:"agent_id"`
	APIKey  string    `json:"api_key"`
}
