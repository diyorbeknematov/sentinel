package analyzer

import (
	"strings"
	"time"

	"github.com/diyorbek/sentinel/internal/models"
)

func (la *LogAnalyzer) AnalyzeNginxLog(log *models.NginxLog) *AnalyzeRes {
	if r := la.checkDDoS(log); r != nil {
		return r
	}
	if r := la.checkDangerousPath(log); r != nil {
		return r
	}
	if r := la.check404Scanning(log); r != nil {
		return r
	}
	return nil
}

// 1. DDoS
// 1 daqiqada 1 IP dan 100+ request
func (la *LogAnalyzer) checkDDoS(log *models.NginxLog) *AnalyzeRes {
	la.mu.Lock()
	defer la.mu.Unlock()

	now := time.Now()

	// request qo‘shamiz
	la.requestCounts[log.IPAddress] = append(la.requestCounts[log.IPAddress], now)
	la.requestCounts[log.IPAddress] = filterRecent(la.requestCounts[log.IPAddress], time.Minute)

	if len(la.requestCounts[log.IPAddress]) >= 100 {

		lastAlert, exists := la.alertedIPs[log.IPAddress]

		// agar oldin alert bo‘lmagan yoki 1 minut o‘tgan bo‘lsa
		if !exists || time.Since(lastAlert) > time.Minute {

			la.alertedIPs[log.IPAddress] = now

			return &AnalyzeRes{
				ThreatType: "DDOS",
				Severity:   "CRITICAL",
				Message:    "IP: " + log.IPAddress + " 1 daqiqada 100+ request",
			}
		}
	}

	return nil
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

func (la *LogAnalyzer) checkDangerousPath(log *models.NginxLog) *AnalyzeRes {
	pathLower := strings.ToLower(log.Path)
	for _, p := range dangerousPaths {
		if strings.HasPrefix(pathLower, p) {
			return &AnalyzeRes{
				ThreatType: "DANGEROUS_PATH",
				Severity:   "HIGH",
				Message:    "IP: " + log.IPAddress + " xavfli sahifaga kirdi: " + log.Path,
			}
		}
	}
	return nil
}

// 3. 404 SCANNING
// 1 IP dan 1 daqiqada 20+ 404
func (la *LogAnalyzer) check404Scanning(log *models.NginxLog) *AnalyzeRes {
	if log.Status != 404 {
		return nil
	}

	la.mu.Lock()
	defer la.mu.Unlock()

	now := time.Now()
	la.notFoundCounts[log.IPAddress] = append(la.notFoundCounts[log.IPAddress], now)
	la.notFoundCounts[log.IPAddress] = filterRecent(la.notFoundCounts[log.IPAddress], time.Minute)

	if len(la.notFoundCounts[log.IPAddress]) >= 20 {
		return &AnalyzeRes{
			ThreatType: "SCANNING",
			Severity:   "WARNING",
			Message:    "IP: " + log.IPAddress + " tizimni paypaslamoqda",
		}
	}
	return nil
}
