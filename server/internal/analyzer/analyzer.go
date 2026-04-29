package analyzer

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AnalyzeRes struct {
	ThreatType string
	Severity   string
	Message    string
}

type LogAnalyzer struct {
	mu sync.Mutex // barcha map'lar uchun bitta mutex yetarli

	// AppLog uchun
	eventCounts map[string][]time.Time
	errorCounts  map[uuid.UUID][]time.Time

	// NginxLog uchun
	requestCounts  map[string][]time.Time
	notFoundCounts map[string][]time.Time
}

func NewLogAnalyzer() *LogAnalyzer {
	return &LogAnalyzer{
		eventCounts:   make(map[string][]time.Time),
		errorCounts:    make(map[uuid.UUID][]time.Time),
		requestCounts:  make(map[string][]time.Time),
		notFoundCounts: make(map[string][]time.Time),
	}
}

func (la *LogAnalyzer) StartCleanup(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			la.mu.Lock()
			// bo'sh bo'lgan keylarni o'chiramiz
			for k, v := range la.eventCounts {
				if len(filterRecent(v, time.Minute)) == 0 {
					delete(la.eventCounts, k)
				}
			}
			for k, v := range la.requestCounts {
				if len(filterRecent(v, time.Minute)) == 0 {
					delete(la.requestCounts, k)
				}
			}
			la.mu.Unlock()
		}
	}
}

func filterRecent(times []time.Time, window time.Duration) []time.Time {
	now := time.Now()
	var recent []time.Time
	for _, t := range times {
		if now.Sub(t) <= window {
			recent = append(recent, t)
		}
	}
	return recent
}
