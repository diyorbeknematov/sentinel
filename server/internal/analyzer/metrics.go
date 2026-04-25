package analyzer

import (
	"fmt"

	"github.com/diyorbek/sentinel/internal/models"
)

func (la *LogAnalyzer) AnalyzeMetric(metric *models.Metric) []*AnalyzeRes {
	var results []*AnalyzeRes

	// Hammasi bir vaqtda xavfli bo'lishi mumkin
	// Shuning uchun slice qaytaramiz
	if r := checkCPU(metric); r != nil {
		results = append(results, r)
	}
	if r := checkRAM(metric); r != nil {
		results = append(results, r)
	}
	if r := checkDisk(metric); r != nil {
		results = append(results, r)
	}

	return results
}

func  checkCPU(metric *models.Metric) *AnalyzeRes {
	if metric.CPU > 90 {
		return &AnalyzeRes{
			ThreatType: "HIGH_CPU",
			Severity:   "HIGH",
			Message:    fmt.Sprintf("CPU yuqori: %.1f%%", metric.CPU),
		}
	}
	return nil
}

func checkRAM(metric *models.Metric) *AnalyzeRes {
	if metric.RAM > 90 {
		return &AnalyzeRes{
			ThreatType: "HIGH_RAM",
			Severity:   "HIGH",
			Message:    fmt.Sprintf("RAM yuqori: %.1f%%", metric.RAM),
		}
	}
	return nil
}

func checkDisk(metric *models.Metric) *AnalyzeRes {
	if metric.Disk > 90 {
		return &AnalyzeRes{
			ThreatType: "HIGH_DISK",
			Severity:   "CRITICAL",
			Message:    fmt.Sprintf("Disk to'ldi: %.1f%%", metric.Disk),
		}
	}
	return nil
}
