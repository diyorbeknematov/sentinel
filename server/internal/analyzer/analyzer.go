package analyzer

import (
	"log/slog"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/service"
)

type LogAnalyzer struct {
	Service *service.Service
	logger  *slog.Logger
}

type AnalyzeRes struct {
	ThreatType string
	Severity   string
	Message    string
}

func NewLogAnalyzer(service *service.Service, logger *slog.Logger) *LogAnalyzer {
	return &LogAnalyzer{
		Service: service,
		logger:  logger,
	}
}

func (la *LogAnalyzer) ProcessAppLog(log *models.Log) {
	res := analyzeAppLog(log)

	// 2. AppLog ga saqlash
	var threatType *string
	if res != nil {
		threatType = &res.ThreatType
	}

	_, err := la.Service.AppLog.CreateAppLog(models.CreateAppLog{
		AgentId:   log.AgentId,
		UserId:    log.UserId,
		Type:      *threatType,
		Level:     log.Level,
		Message:   log.Message,
		IPAddress: log.IPAddress,
	})

	if err != nil {
		la.logger.Error(err.Error())
		panic(err)
	}

	if res != nil {
		_, err = la.Service.Alert.CreateAlert(models.CreateAlert{
			AgentId:  log.AgentId,
			Type:     res.ThreatType,
			Message:  res.Message,
			Severity: res.Severity,
		})
		if err != nil {
			la.logger.Error(err.Error())
			panic(err)
		}
	}
}

func (la *LogAnalyzer) ProcessNginxLog(log *models.NginxLog) {
	_ = analyzeNginxLog(log)

	// 2. nginxlogs ga saqlash
	// var threatType *string
	// if res != nil {
	// 	threatType = &res.ThreatType
	// }

}

func (la *LogAnalyzer) ProcessMetric(metric *models.Metric) {
	_, err := la.Service.Metric.CreateMetric(models.CreateMetric{
		AgentId: metric.AgentId,
		CPU:     metric.CPU,
		RAM:     metric.RAM,
		Disk:    metric.Disk,
	})

	if err != nil {
		la.logger.Error(err.Error())
		panic(err)
	}

	results := analyzeMetric(metric)
	for _, res := range results {
		_, err = la.Service.Alert.CreateAlert(models.CreateAlert{
			AgentId:  metric.AgentId,
			Type:     res.ThreatType,
			Message:  res.Message,
			Severity: res.Severity,
		})

		if err != nil {
			la.logger.Error(err.Error())
			panic(err)
		}
	}
}
