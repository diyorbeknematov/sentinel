package analyzer

import (
	"strings"
	"time"

	"github.com/diyorbek/sentinel/internal/models"
)

var sqlPatterns = []string{
	"' or ",
	"' and ",
	"union select",
	"drop table",
	"insert into",
	"delete from",
	"'; --",
	"1=1",
	"exec(",
	"xp_cmdshell",
	"/*",
}

func (la *LogAnalyzer) AnalyzeAppLog(log *models.Log) *AnalyzeRes {
	if r := la.checkBruteForce(log); r != nil {
		return r
	}
	if r := la.checkSQLInjection(log); r != nil {
		return r
	}
	if r := la.checkHighErrorRate(log); r != nil {
		return r
	}
	return nil
}

// 1. BRUTE FORCE
// 1 daqiqada 10+ login_failed — bir IP dan
func (la *LogAnalyzer) checkBruteForce(log *models.Log) *AnalyzeRes {
	if log.Event != "login_failed" {
		return nil
	}

	la.mu.Lock()
	defer la.mu.Unlock()

	now := time.Now()
	la.failedLogins[log.UserId] = append(la.failedLogins[log.UserId], now)
	la.failedLogins[log.UserId] = filterRecent(la.failedLogins[log.UserId], time.Minute)

	if len(la.failedLogins[log.UserId]) >= 10 {
		return &AnalyzeRes{
			ThreatType: "BRUTE_FORCE",
			Severity:   "CRITICAL",
			Message:    "user id: " + log.UserId + " 1 daqiqada 10+ marta login urinish",
		}
	}
	return nil

}

// 2. SQL INJECTION
// Message da xavfli pattern bormi
func (la *LogAnalyzer) checkSQLInjection(log *models.Log) *AnalyzeRes {
	msgLower := strings.ToLower(log.Message)
	for _, pattern := range sqlPatterns {
		if strings.Contains(msgLower, pattern) {
			return &AnalyzeRes{
				ThreatType: "SQL_INJECTION",
				Severity:   "CRITICAL",
				Message:    "user id: " + log.UserId + " xavfli so'rov: " + log.Message,
			}
		}
	}
	return nil
}

// 3. KO'P ERROR
// 5 daqiqada 10+ ERROR — bir agent dan

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
			Severity:   "MEDIUM",
			Message:    "Agent ko'p xato: 5 daqiqada 10+ ERROR",
		}
	}
	return nil
}
