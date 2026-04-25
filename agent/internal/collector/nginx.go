package collector

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/diyorbek/sentinel/agent/internal/config"
	"github.com/diyorbek/sentinel/agent/internal/models"
)

var nginxRegex = regexp.MustCompile(
	`(?P<ip>\S+) - - \[(?P<time>.*?)\] "(?P<method>\S+) (?P<path>\S+) \S+" (?P<status>\d+) (?P<bytes>\d+) "[^"]*" "(?P<ua>.*?)"`,
)

func StartNginxLogCollector(ctx context.Context, cfg *config.Config, eventCh chan<- models.Event) {
	file, err := os.Open(cfg.NginxLog)
	if err != nil {
		slog.Error("open nginx log failed", "path", cfg.NginxLog, "err", err)
		return
	}
	defer file.Close()

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		slog.Error("seek failed", "err", err)
		return
	}

	slog.Info("watching nginx log", "path", cfg.NginxLog)

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if len(line) > 0 {
			nginxLog, parseErr := parseNginxLog(line)
			if parseErr != nil {
				slog.Warn("failed to parse nginx log", "err", parseErr)
			} else {
				// channel'ga yuboramiz
				select {
				case eventCh <- models.Event{
					Type:    models.EventNginxLog,
					AgentID: cfg.AgentID,
					Payload: nginxLog,
				}:
				case <-ctx.Done():
					return
				}
			}
		}

		if err != nil {
			if err == io.EOF {
				// yangi satr yo'q — biroz kutib qayta urinamiz
				select {
				case <-time.After(200 * time.Millisecond):
					continue
				case <-ctx.Done():
					return
				}
			}
			slog.Error("read nginx log failed", "err", err)
			return
		}
	}
}

func parseNginxLog(line string) (*models.NginxLogPayload, error) {
	match := nginxRegex.FindStringSubmatch(line)
	if match == nil {
		return nil, fmt.Errorf("line does not match nginx log format: %q", line)
	}

	groups := make(map[string]string)
	for i, name := range nginxRegex.SubexpNames() {
		if i != 0 && name != "" {
			groups[name] = match[i]
		}
	}

	status, err := strconv.Atoi(groups["status"])
	if err != nil {
		return nil, fmt.Errorf("invalid status %q: %w", groups["status"], err)
	}

	bodyBytes, err := strconv.Atoi(groups["bytes"])
	if err != nil {
		return nil, fmt.Errorf("invalid bytes %q: %w", groups["bytes"], err)
	}

	logTime, err := time.Parse("02/Jan/2006:15:04:05 -0700", groups["time"])
	if err != nil {
		return nil, fmt.Errorf("invalid time %q: %w", groups["time"], err)
	}

	return &models.NginxLogPayload{
		IP:        groups["ip"],
		Method:    groups["method"],
		Path:      groups["path"],
		Status:    status,
		Bytes:     bodyBytes,
		UserAgent: groups["ua"],
		LogTime:   logTime,
	}, nil
}
