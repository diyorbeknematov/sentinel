package service

import (
	"time"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/pkg/apperrors"
	"github.com/google/uuid"
)

type StatsService struct {
	repos *repository.Repository
	cfg   *config.Config
}

func NewStatsService(repos *repository.Repository, cfg *config.Config) *StatsService {
	return &StatsService{
		repos: repos,
		cfg:   cfg,
	}
}

func (s *StatsService) GetDashboardStats(filter models.StatsFilter) (models.DashboardStats, error) {
	// Period setup
	duration := time.Hour
	switch filter.Period {
	case "24h":
		duration = 24 * time.Hour
	case "7d":
		duration = 7 * 24 * time.Hour
	}
	
	return s.repos.Stats.GetDashboardStats(filter, duration)
}

func (s *StatsService) GetLogVolumeStats(id string) ([]models.LogVolumeStats, error) {
	var agentID uuid.UUID = uuid.Nil
	if id != "" {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			return nil, apperrors.BadRequest("noto'g'ri UUID format")
		}
		agentID = parsedID
	}

	logVolumes, err := s.repos.Stats.GetLogVolume(agentID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return logVolumes, nil
}
