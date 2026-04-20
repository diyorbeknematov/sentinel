package analyzer

import (
	"strings"
	"time"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/google/uuid"
)

var (
	// Brute Force: IP → [vaqtlar]
	failedLogins = make(map[string][]time.Time)

	// Ko'p ERROR: agentID → [vaqtlar]
	errorCounts = make(map[uuid.UUID][]time.Time)
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

func analyzeAppLog(log *models.Log) *AnalyzeRes {
	// 1. Brute Force
	if checkBruteForce(log) {
		return &AnalyzeRes{
			ThreatType: "BRUTE_FORCE",
			Severity:   "CRITICAL",
			Message:    "user id: " + log.UserId + " 1 daqiqada 10+ marta login urinish",
		}
	}

	// 2. SQL Injection
	if checkSQLInjection(log) {
		return &AnalyzeRes{
			ThreatType: "SQL_INJECTION",
			Severity:   "CRITICAL",
			Message:    "user id: " + log.UserId + " xavfli so'rov: " + log.Message,
		}
	}

	// 3. Ko'p ERROR
	if checkHighErrorRate(log) {
		return &AnalyzeRes{
			ThreatType: "HIGH_ERROR_RATE",
			Severity:   "MEDIUM",
			Message:    "Agent ko'p xato: 1 daqiqada 10+ ERROR",
		}
	}

	return nil
}

// 1. BRUTE FORCE
// 1 daqiqada 10+ login_failed — bir IP dan
func checkBruteForce(log *models.Log) bool {
	if log.Event != "login_failed" {
		return false
	}

	now := time.Now()
	userId := log.UserId

	failedLogins[userId] = append(failedLogins[userId], now)

	var recent []time.Time
	for _, t := range failedLogins[userId] {
		if now.Sub(t) <= time.Minute {
			recent = append(recent, t)
		}
	}
	failedLogins[userId] = recent

	// 10+ bo'lsa Brute Force
	return len(failedLogins[userId]) >= 10
}

// 2. SQL INJECTION
// Message da xavfli pattern bormi
func checkSQLInjection(log *models.Log) bool {
	msgLower := strings.ToLower(log.Message)
	for _, pattern := range sqlPatterns {
		if strings.Contains(msgLower, pattern) {
			return true
		}
	}
	return false
}

// 3. KO'P ERROR
// 5 daqiqada 10+ ERROR — bir agent dan

func checkHighErrorRate(log *models.Log) bool {
	if log.Level != "ERROR" {
		return false
	}

	now := time.Now()
	id := log.AgentId

	errorCounts[id] = append(errorCounts[id], now)

	var recent []time.Time
	for _, t := range errorCounts[id] {
		if now.Sub(t) <= 5*time.Minute {
			recent = append(recent, t)
		}
	}
	errorCounts[id] = recent

	return len(errorCounts[id]) >= 10
}
