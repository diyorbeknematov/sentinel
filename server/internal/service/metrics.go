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

type metricService struct {
	repo   *repository.Repository
	cfg    *config.Config
	logger *slog.Logger
}

func NewMetricService(repo *repository.Repository, cfg *config.Config, logger *slog.Logger) *metricService {
	return &metricService{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *metricService) CreateMetric(metric models.CreateMetric) (uuid.UUID, error) {
	id, err := s.repo.Metric.CreateMetric(metric)
	if err != nil {
		return uuid.Nil, apperrors.Internal(err)
	}

	return id, nil
}

func (s *metricService) GetMetricsByID(id uuid.UUID) (models.Metric, error) {
	metric, err := s.repo.Metric.GetMetricsByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Metric{}, apperrors.NotFound("Metric topilmadi", err)
		}
		return models.Metric{}, apperrors.Internal(err)
	}

	return metric, nil
}

func (s *metricService) ListMetrics(filter models.FilterMetrics) ([]models.Metric, int, error) {
	metrics, total, err := s.repo.Metric.ListMetrics(filter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Metric{}, 0, apperrors.NotFound("Metric topilmadi", err)
		}
		return []models.Metric{}, 0, apperrors.Internal(err)
	}

	return metrics, total, nil
}
