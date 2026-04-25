package collector

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"time"

	"github.com/diyorbek/sentinel/agent/internal/config"
	"github.com/diyorbek/sentinel/agent/internal/models"
)

var re = regexp.MustCompile(`(?P<time>\S+) (?P<level>\S+) (?P<event>\S+) (?P<message>.*)`)

func StartAppLogCollector(ctx context.Context, cfg *config.Config, eventCh chan<- models.Event) {
	file, err := os.Open(cfg.AppLog)
	if err != nil {
		slog.Error("open app log failed", "path", cfg.AppLog, "err", err)
		return
	}
	defer file.Close()

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		slog.Error("seek failed", "err", err)
		return
	}

	slog.Info("watching app log", "path", cfg.AppLog)

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if len(line) > 0 {
			appLog, parseErr := parseAppLog(line)
			if parseErr != nil {
				slog.Warn("failed to parse log line", "err", parseErr)
			} else {
				// TODO: send log
				select {
				case eventCh <- models.Event{
					Type:    models.EventAppLog,
					AgentID: cfg.AgentID,
					Payload: appLog,
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
			slog.Error("read app log failed", "err", err)
			return
		}
	}
}

func parseAppLog(line string) (*models.AppLogPayload, error) {
	if isJSON(line) {
		return parseJSON(line)
	}
	return parseText(line)
}

func isJSON(line string) bool {
	return json.Valid([]byte(line))
}

func parseJSON(line string) (*models.AppLogPayload, error) {
	var payload models.AppLogPayload
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		return nil, fmt.Errorf("json parse: %w", err)
	}
	return &payload, nil
}

func parseText(line string) (*models.AppLogPayload, error) {
	match := re.FindStringSubmatch(line)
	if match == nil {
		return nil, fmt.Errorf("line does not match expected format: %q", line)
	}

	groups := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && name != "" {
			groups[name] = match[i]
		}
	}

	logTime, err := time.Parse("02/Jan/2006:15:04:05 -0700", groups["time"])
	if err != nil {
		return nil, fmt.Errorf("parse log time %q: %w", groups["time"], err)
	}

	return &models.AppLogPayload{
		Event:   groups["event"],
		Level:   groups["level"],
		Message: groups["message"],
		LogTime: logTime,
	}, nil
}
