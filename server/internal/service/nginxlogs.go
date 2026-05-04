package service

import (
	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/pkg/apperrors"
	"github.com/google/uuid"
)

type nginxLogService struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewNginxLogService(repo *repository.Repository, cfg *config.Config) *nginxLogService {
	return &nginxLogService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *nginxLogService) CreateNginxLog(nlog models.CreateNginxLog) (uuid.UUID, error) {
	nlog.Id = uuid.New()
	// 1. DB ga yozish
	id, err := s.repo.NginxLog.CreateNginxLog(nlog)
	if err != nil {
		return uuid.Nil, apperrors.Internal(err)
	}

	return id, nil
}

func (s *nginxLogService) GetNginxLogByID(id string) (models.NginxLog, error) {
	UId, err := uuid.Parse(id)
	if err != nil {
		return models.NginxLog{}, apperrors.BadRequest("noto'g'ri UUID format")
	}
	// 1. DB dan olish
	nlog, err := s.repo.NginxLog.GetNginxLogByID(UId)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNoRowsAffected) {
			return models.NginxLog{}, apperrors.NotFound("nginxlog topilmadi", err)
		}
		return models.NginxLog{}, apperrors.Internal(err)
	}

	return nlog, nil
}

func (s *nginxLogService) ListNginxLogs(filter models.FilterNginxLog) ([]models.NginxLogResponse, int, error) {
	agentUId := uuid.Nil
	if filter.AgentId != "" {
		agentId, err := uuid.Parse(filter.AgentId)
		if err != nil {
			return []models.NginxLogResponse{}, 0, apperrors.BadRequest("noto'g'ri UUID format")
		}
		agentUId = agentId
	}

	// 1. DB dan olish
	nlogs, total, err := s.repo.NginxLog.ListNginxLogs(models.FilterNginxLogDB{
		AgentId: agentUId,
		Method:  filter.Method,
		Status:  filter.Status,
		From:    filter.From,
		To:      filter.To,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
	})
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNoRowsAffected) {
			return []models.NginxLogResponse{}, 0, apperrors.NotFound("nginxlogs topilmadi", err)
		}
		return []models.NginxLogResponse{}, 0, apperrors.Internal(err)
	}

	return nlogs, total, nil
}
