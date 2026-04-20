package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/diyorbek/sentinel/agent/internal/models"
)

func Register(serverURL string, req models.RegisterRequest) (models.RegisterResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return models.RegisterResponse{}, err
	}

	resp, err := http.Post(serverURL+"/sentinel/api/agents", "application/json", bytes.NewBuffer(body))
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
