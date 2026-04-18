package service

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/pkg/apperrors"
	"github.com/google/uuid"
)

type agentService struct {
	repo   *repository.Repository
	logger *slog.Logger
	cfg    *config.Config
}

func NewAgentService(repo *repository.Repository, cfg *config.Config, logger *slog.Logger) *agentService {
	return &agentService{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *agentService) CreateAgent(req models.CreateAgent) (uuid.UUID, error) {
	apiKey, err := generateAPIKey()
	if err != nil {
		return uuid.Nil, apperrors.Internal(err)
	}

	req.APIKey = apiKey

	id, err := s.repo.CreateAgent(req)
	if err != nil {
		return uuid.Nil, apperrors.Internal(err)
	}

	return id, nil
}

func (s *agentService) GetAgentByID(id uuid.UUID) (models.Agent, error) {
	agent, err := s.repo.Agent.GetAgentByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Agent{}, apperrors.NotFound(
				"agent topilmadi",
				err,
			)
		}
		return models.Agent{}, apperrors.Internal(err)
	}

	return agent, nil
}

func (s *agentService) GetAgentByAPIKey(apiKey string) (models.Agent, error) {
	agent, err := s.repo.Agent.GetAgentByAPIKey(apiKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Agent{}, apperrors.NotFound("agent topilmadi", err)
		}
		return models.Agent{}, apperrors.Internal(err)
	}

	return agent, nil
}

func (s *agentService) ListAgents(filter models.FilterAgent) ([]models.Agent, int, error) {
	agents, total, err := s.repo.Agent.ListAgents(filter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Agent{}, 0, apperrors.NotFound("agent topilmadi", err)
		}
		return []models.Agent{}, 0, apperrors.Internal(err)
	}

	return agents, total, nil
}

func (s *agentService) UpdateLastSeen(lastSeen models.UpdateLastSeen) error {
	err := s.repo.Agent.UpdateLastSeen(lastSeen)
	if err != nil {
		if apperrors.Is(err, apperrors.EerrNoRowsAffected) {
			return apperrors.BadRequest("bunday ma'lumotga ega agent yo'q")
		}
		return apperrors.Internal(err)
	}

	return nil
}

func (s *agentService) DeleteAgent(id uuid.UUID) (err error) {
	err = s.repo.Agent.DeleteAgent(id)
	if err != nil {
		if apperrors.Is(err, apperrors.EerrNoRowsAffected) {
			return apperrors.BadRequest("bunday agent yo'q")
		}

		return apperrors.Internal(err)
	}

	return
}
