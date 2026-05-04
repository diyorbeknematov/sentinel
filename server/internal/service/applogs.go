package service

import (
	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/pkg/apperrors"
	"github.com/google/uuid"
)

type logService struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewLogService(repo *repository.Repository, cfg *config.Config) *logService {
	return &logService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *logService) CreateAppLog(log models.CreateAppLog) (uuid.UUID, error) {
	// 1. DB ga yozish
	log.Id = uuid.New()
	id, err := s.repo.AppLog.CreateAppLog(log)
	if err != nil {
		return uuid.Nil, apperrors.Internal(err)
	}

	return id, nil
}

func (s *logService) GetLogByID(id string) (models.Log, error) {
	UId, err := uuid.Parse(id)
	if err != nil {
		return models.Log{}, apperrors.BadRequest("noto'g'ri UUID format")
	}
	log, err := s.repo.AppLog.GetLogByID(UId)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNoRowsAffected) {
			return models.Log{}, apperrors.NotFound("log topilmadi", err)
		}
		return models.Log{}, apperrors.Internal(err)
	}

	return log, nil
}

func (s *logService) ListLogs(filter models.FilterAppLog) ([]models.AppLogResponse, int, error) {
	agentUId := uuid.Nil
	if filter.AgentId != "" {
		agentId, err := uuid.Parse(filter.AgentId)
		if err != nil {
			return []models.AppLogResponse{}, 0, apperrors.BadRequest("noto'g'ri UUID format")
		}
		agentUId = agentId
	}

	logs, total, err := s.repo.AppLog.ListLogs(models.FilterAppLogDB{
		AgentId: agentUId,
		UserId:  filter.UserId,
		Event:   filter.Event,
		Level:   filter.Level,
		From:    filter.From,
		To:      filter.To,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
	})
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNoRowsAffected) {
			return []models.AppLogResponse{}, 0, apperrors.NotFound("log topilmadi", err)
		}
		return []models.AppLogResponse{}, 0, apperrors.Internal(err)
	}

	return logs, total, nil
}
