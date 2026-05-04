package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/pkg/apperrors"
	"github.com/google/uuid"
)

type agentService struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewAgentService(repo *repository.Repository, cfg *config.Config) *agentService {
	return &agentService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *agentService) CreateAgent(req models.CreateAgent) (models.CreateAgentResponse, error) {
	UId, err := uuid.Parse(req.Id)
	if err != nil {
		return models.CreateAgentResponse{}, apperrors.BadRequest("noto'g'ri UUID format")
	}

	id, err := s.repo.Agent.CreateAgent(models.CreateAgentDB{
		Id:        uuid.New(),
		AccountID: UId,
		Name:      req.Name,
		IPAddress: req.IPAddress,
		LastSeen:  time.Now(),
	})
	if err != nil {
		return models.CreateAgentResponse{}, apperrors.Internal(err)
	}

	return models.CreateAgentResponse{
		ID:           id,
		KafkaBrokers: s.cfg.KafkaExternalBrokers,
		KafkaTopic:   s.cfg.KafkaTopic,
	}, nil
}

func (s *agentService) GetAgentByID(id string) (models.Agent, error) {
	UId, err := uuid.Parse(id)
	if err != nil {
		return models.Agent{}, apperrors.BadRequest("noto'g'ri UUID format")
	}

	agent, err := s.repo.Agent.GetAgentByID(UId)
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

func (s *agentService) ListAgents(filter models.FilterAgent) ([]models.Agent, int, error) {
	var accountID uuid.UUID = uuid.Nil

	if filter.AccountID != "" {
		id, err := uuid.Parse(filter.AccountID)
		if err != nil {
			return []models.Agent{}, 0, apperrors.BadRequest("noto'g'ri UUID format")
		}
		accountID = id
	}

	agents, total, err := s.repo.Agent.ListAgents(models.FilterAgentDB{
		AccountID: accountID,
		Name:      filter.Name,
		From:      filter.From,
		To:        filter.To,
		Limit:     filter.Limit,
		Offset:    filter.Offset,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Agent{}, 0, apperrors.NotFound("agent topilmadi", err)
		}
		return []models.Agent{}, 0, apperrors.Internal(err)
	}

	now := time.Now().UTC()
	for i := range agents {
		if agents[i].LastSeen.IsZero() {
			agents[i].Status = "offline"
			continue
		}

		if now.Sub(agents[i].LastSeen) < 60*time.Second {
			agents[i].Status = "online"
		} else {
			agents[i].Status = "offline"
		}
	}

	return agents, total, nil
}

func (s *agentService) UpdateLastSeen(agentID uuid.UUID) error {
	err := s.repo.Agent.UpdateLastSeen(models.UpdateLastSeen{
		Id:       agentID,
		LastSeen: time.Now().UTC(),
	})
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNoRowsAffected) {
			return apperrors.BadRequest("bunday ma'lumotga ega agent yo'q")
		}
		return apperrors.Internal(err)
	}

	return nil
}

func (s *agentService) DeleteAgent(id string) (err error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return apperrors.BadRequest("noto'g'ri UUID format")
	}

	err = s.repo.Agent.DeleteAgent(uuid)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNoRowsAffected) {
			return apperrors.BadRequest("bunday agent yo'q")
		}

		return apperrors.Internal(err)
	}

	return
}
