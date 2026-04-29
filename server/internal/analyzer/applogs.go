package analyzer

import (
	"time"

	"github.com/diyorbek/sentinel/internal/models"
)

func (la *LogAnalyzer) AnalyzeAppLog(log *models.Log) *AnalyzeRes {
	if r := la.checkHighErrorRate(log); r != nil {
		return r
	}

	if r := la.checkEventSpike(log); r != nil {
		return r
	}

	return nil
}

func (la *LogAnalyzer) checkHighErrorRate(log *models.Log) *AnalyzeRes {
	if log.Level != "ERROR" {
		return nil
	}

	la.mu.Lock()
	defer la.mu.Unlock()

	now := time.Now()
	la.errorCounts[log.AgentId] = append(la.errorCounts[log.AgentId], now)
	la.errorCounts[log.AgentId] = filterRecent(la.errorCounts[log.AgentId], 5*time.Minute)

	if len(la.errorCounts[log.AgentId]) >= 10 {
		return &AnalyzeRes{
			ThreatType: "HIGH_ERROR_RATE",
			Severity:   "WARNING",
			Message:    "Agent ko'p xato: 5 daqiqada 10+ ERROR",
		}
	}
	return nil
}

func (la *LogAnalyzer) checkEventSpike(log *models.Log) *AnalyzeRes {
	la.mu.Lock()
	defer la.mu.Unlock()

	now := time.Now()

	key := log.Event

	la.eventCounts[key] = append(la.eventCounts[key], now)
	la.eventCounts[key] = filterRecent(la.eventCounts[key], 1*time.Minute)

	if len(la.eventCounts[key]) >= 50 {
		return &AnalyzeRes{
			ThreatType: "EVENT_SPIKE",
			Severity:   "HIGH",
			Message:    "Event spike detected: " + log.Event,
		}
	}

	return nil
}
