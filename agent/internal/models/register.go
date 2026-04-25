package models

import "github.com/google/uuid"

type RegisterRequest struct {
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}

type RegisterResponse struct {
	AgentID uuid.UUID `json:"id"`
	APIKey  string    `json:"api_key"`
}
