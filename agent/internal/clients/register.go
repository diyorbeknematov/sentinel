package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/diyorbek/sentinel/agent/internal/config"
	"github.com/diyorbek/sentinel/agent/internal/models"
)

func Register(serverURL, apiKey string, req models.RegisterRequest) (models.RegisterResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return models.RegisterResponse{}, err
	}

	reqHTTP, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/sentinel/agents",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return models.RegisterResponse{}, err
	}

	// AAD API KEY HEADER
	reqHTTP.Header.Set("Content-Type", "application/json")
	reqHTTP.Header.Set("X-API-Key", apiKey) 

	client := &http.Client{}
	resp, err := client.Do(reqHTTP)
	if err != nil {
		return models.RegisterResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return models.RegisterResponse{}, fmt.Errorf("register failed: status %d", resp.StatusCode)
	}

	var res models.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return models.RegisterResponse{}, err
	}

	return res, nil
}

func SendHeartbeat(cfg *config.Config) error {
	body, err := json.Marshal(models.HeartbeatRequest{
		AgentID: cfg.AgentID,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		cfg.ServerURL+"/sentinel/agents/heartbeat",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", cfg.APIKey)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("heartbeat failed: status %d", resp.StatusCode)
	}

	return nil
}
