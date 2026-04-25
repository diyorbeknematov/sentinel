package collector

import (
    "context"
    "fmt"
    "log/slog"
    "time"

    "github.com/diyorbek/sentinel/agent/internal/config"
    "github.com/diyorbek/sentinel/agent/internal/models"
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/disk"
    "github.com/shirou/gopsutil/v3/mem"
)

func StartMetricsCollector(ctx context.Context, cfg *config.Config, eventCh chan<- models.Event) {
    slog.Info("metrics collector started")

    ticker := time.NewTicker(cfg.MetricsInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            slog.Info("metrics collector stopped")
            return

        case <-ticker.C:
            metric, err := collectMetrics()
            if err != nil {
                slog.Warn("collect metrics failed", "err", err)
                continue
            }

            select {
            case eventCh <- models.Event{
                Type:    models.EventMetric,
                AgentID: cfg.AgentID,
                Payload: metric,
            }:
            case <-ctx.Done():
                return
            }
        }
    }
}

func collectMetrics() (*models.MetricPayload, error) {
    cpuPercent, err := cpu.Percent(time.Second, false)
    if err != nil {
        return nil, fmt.Errorf("cpu: %w", err)
    }
    if len(cpuPercent) == 0 {
        return nil, fmt.Errorf("cpu: empty result")
    }

    vm, err := mem.VirtualMemory()
    if err != nil {
        return nil, fmt.Errorf("ram: %w", err)
    }

    diskStat, err := disk.Usage("/")
    if err != nil {
        return nil, fmt.Errorf("disk: %w", err)
    }

    return &models.MetricPayload{
        CPU:     cpuPercent[0],
        RAM:     vm.UsedPercent,
        Disk:    diskStat.UsedPercent,
        LogTime: time.Now(),
    }, nil
}