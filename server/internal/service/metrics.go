package service

import (
	"database/sql"
	"errors"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/pkg/apperrors"
	"github.com/google/uuid"
)

type metricService struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewMetricService(repo *repository.Repository, cfg *config.Config) *metricService {
	return &metricService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *metricService) CreateMetric(metric models.CreateMetric) (uuid.UUID, error) {
	metric.Id = uuid.New()
	id, err := s.repo.Metric.CreateMetric(metric)
	if err != nil {
		return uuid.Nil, apperrors.Internal(err)
	}

	return id, nil
}

func (s *metricService) GetMetricsByID(id string) (models.Metric, error) {
	UId, err := uuid.Parse(id)
	if err != nil {
		return models.Metric{}, apperrors.BadRequest("noto'g'ri UUID format")
	}

	metric, err := s.repo.Metric.GetMetricsByID(UId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Metric{}, apperrors.NotFound("Metric topilmadi", err)
		}
		return models.Metric{}, apperrors.Internal(err)
	}

	return metric, nil
}

func (s *metricService) ListMetrics(filter models.FilterMetrics) ([]models.MetricResponse, int, error) {
	agentUId := uuid.Nil
	if filter.AgentId != "" {
		agentId, err := uuid.Parse(filter.AgentId)
		if err != nil {
			return []models.MetricResponse{}, 0, apperrors.BadRequest("noto'g'ri UUID format")
		}
		agentUId = agentId
	}

	metrics, total, err := s.repo.Metric.ListMetrics(models.FilterMetricsDB{
		AgentId: agentUId,
		From:    filter.From,
		To:      filter.To,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.MetricResponse{}, 0, apperrors.NotFound("Metric topilmadi", err)
		}
		return []models.MetricResponse{}, 0, apperrors.Internal(err)
	}

	return metrics, total, nil
}
