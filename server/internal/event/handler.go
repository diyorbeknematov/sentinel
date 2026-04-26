package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/diyorbek/sentinel/internal/analyzer"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/service"
)

type EventHandler struct {
	service  *service.Service
	analyzer *analyzer.LogAnalyzer
	logger   *slog.Logger
}

func NewEventHanler(service *service.Service, analyzer *analyzer.LogAnalyzer, logger *slog.Logger) *EventHandler {
	return &EventHandler{
		service:  service,
		analyzer: analyzer,
		logger:   logger,
	}
}

func (h *EventHandler) HandleMetric(ctx context.Context, event models.Event) error {
	fmt.Println(event)
	var metric models.Metric
	if err := json.Unmarshal(event.Payload, &metric); err != nil {
		return fmt.Errorf("parse metric: %w", err)
	}

	// 1. DB ga yozish
	if _, err := h.service.Metric.CreateMetric(models.CreateMetric{
		AgentId: event.AgentID,
		CPU:     metric.CPU,
		RAM:     metric.RAM,
		Disk:    metric.Disk,
		LogTime: metric.LogTime,
	}); err != nil {
		return fmt.Errorf("save metric: %w", err)
	}

	// 2. Analyze qilish
	results := h.analyzer.AnalyzeMetric(&metric)
	for _, res := range results {
		if _, err := h.service.Alert.CreateAlert(models.CreateAlert{
			AgentId:  event.AgentID,
			Type:     res.ThreatType,
			Message:  res.Message,
			Severity: res.Severity,
		}); err != nil {
			h.logger.Warn("analyze failed", "err", err)
		}
	}

	return nil
}
func (h *EventHandler) HandleNginxLog(ctx context.Context, event models.Event) error {
	var nginxLog models.NginxLog

	if err := json.Unmarshal(event.Payload, &nginxLog); err != nil {
		return fmt.Errorf("parse nginx log: %w", err)
	}

	nginxLog.AgentId = event.AgentID

	if _, err := h.service.NginxLog.CreateNginxLog(models.CreateNginxLog{
		AgentId:   nginxLog.AgentId,
		IPAddress: nginxLog.IPAddress,
		Method:    nginxLog.Method,
		Path:      nginxLog.Path,
		Status:    nginxLog.Status,
		Bytes:     nginxLog.Bytes,
		UserAgent: nginxLog.UserAgent,
		LogTime:   nginxLog.LogTime,
	}); err != nil {
		return fmt.Errorf("save nginx log: %w", err)
	}

	res := h.analyzer.AnalyzeNginxLog(&nginxLog)
	if res != nil {
		if _, err := h.service.Alert.CreateAlert(models.CreateAlert{
			AgentId:  nginxLog.AgentId,
			Type:     res.ThreatType,
			Message:  res.Message,
			Severity: res.Severity,
		}); err != nil {
			h.logger.Warn("analyze failed", "err", err)
		}
	}

	return nil
}

func (h *EventHandler) HandleAppLog(ctx context.Context, event models.Event) error {
	var appLog models.Log

	if err := json.Unmarshal(event.Payload, &appLog); err != nil {
		return fmt.Errorf("parse app log: %w", err)
	}

	if _, err := h.service.AppLog.CreateAppLog(models.CreateAppLog{
		AgentId: event.AgentID,
		UserId:  appLog.UserId,
		Event:   appLog.Event,
		Level:   appLog.Level,
		Message: appLog.Message,
		LogTime: appLog.LogTime,
	}); err != nil {
		return fmt.Errorf("save app log: %w", err)
	}

	res := h.analyzer.AnalyzeAppLog(&appLog)
	if res != nil {
		if _, err := h.service.CreateAlert(models.CreateAlert{
			AgentId:  event.AgentID,
			Type:     res.ThreatType,
			Message:  res.Message,
			Severity: res.Severity,
		}); err != nil {
			h.logger.Warn("analyze failed", "err", err)
		}
	}

	return nil
}
