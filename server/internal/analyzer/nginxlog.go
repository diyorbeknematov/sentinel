package analyzer

import (
	"strings"
	"time"

	"github.com/diyorbek/sentinel/internal/models"
)

var (
	requestCounts  = make(map[string][]time.Time) // ip → vaqtlar (DDoS)
	notFoundCounts = make(map[string][]time.Time) // ip → vaqtlar (404)
)

func analyzeNginxLog(log *models.NginxLog) *AnalyzeRes {

	// 1. DDoS
	if checkDDoS(log) {
		return &AnalyzeRes{
			ThreatType: "DDOS",
			Severity:   "CRITICAL",
			Message:    "IP: " + log.IPAddress + " 1 daqiqada 100+ request",
		}
	}

	// 2. Xavfli URL
	if checkDangerousPath(log) {
		return &AnalyzeRes{
			ThreatType: "DANGEROUS_PATH",
			Severity:   "HIGH",
			Message:    "IP: " + log.IPAddress + " xavfli sahifaga kirdi: " + log.Path,
		}
	}

	// 3. 404 Scanning
	if check404Scanning(log) {
		return &AnalyzeRes{
			ThreatType: "SCANNING",
			Severity:   "MEDIUM",
			Message:    "IP: " + log.IPAddress + " tizimni paypaslamoqda",
		}
	}

	return nil
}

// 1. DDoS
// 1 daqiqada 1 IP dan 100+ request
func checkDDoS(log *models.NginxLog) bool {
	now := time.Now()
	ip := log.IPAddress

	requestCounts[ip] = append(requestCounts[ip], now)

	var recent []time.Time
	for _, t := range requestCounts[ip] {
		if now.Sub(t) <= 1*time.Minute {
			recent = append(recent, t)
		}
	}
	requestCounts[ip] = recent

	return len(requestCounts[ip]) >= 100
}

// 2. XAVFLI URL
// /admin, /.env kabi sahifalarga kirish urinishi
var dangerousPaths = []string{
	"/admin",
	"/.env",
	"/config",
	"/wp-admin",
	"/phpmyadmin",
	"/.git",
}

func checkDangerousPath(log *models.NginxLog) bool {
	pathLower := strings.ToLower(log.Path)
	for _, p := range dangerousPaths {
		if strings.HasPrefix(pathLower, p) {
			return true
		}
	}
	return false
}

// 3. 404 SCANNING
// 1 IP dan 1 daqiqada 20+ 404
func check404Scanning(log *models.NginxLog) bool {
	if log.Status != 404 {
		return false
	}

	now := time.Now()
	ip := log.IPAddress

	notFoundCounts[ip] = append(notFoundCounts[ip], now)

	var recent []time.Time
	for _, t := range notFoundCounts[ip] {
		if now.Sub(t) <= 1*time.Minute {
			recent = append(recent, t)
		}
	}
	notFoundCounts[ip] = recent

	return len(notFoundCounts[ip]) >= 20
}
