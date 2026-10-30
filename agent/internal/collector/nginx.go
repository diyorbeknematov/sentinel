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
	`(?P<ip>\S+) - - \[(?P<time>[^\]]+)\] "(?P<method>\S+) (?P<path>\S+) (?P<protocol>\S+)" (?P<status>\d{3}) (?P<bytes>\d+|-) "(?P<referrer>[^"]*)" "(?P<ua>[^"]*)"`,
)

func StartNginxLogCollector(ctx context.Context, cfg *config.Config, eventCh chan<- models.Event) {
	file, err := os.Open(cfg.NginxLog)
	if err != nil {
		slog.Error("open nginx log failed", "path", cfg.NginxLog, "err", err)
		return
	}
	// file is reopened when the log is rotated, so close whichever one is open at the end.
	defer func() { file.Close() }()

	// offset is how many bytes of the file have been read; it is used to notice a truncated file.
	offset, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		slog.Error("seek failed", "err", err)
		return
	}

	slog.Info("watching nginx log", "path", cfg.NginxLog)

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		offset += int64(len(line))

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
				// logrotate faylni almashtirgan yoki tozalagan bo'lsa, yangisini boshidan ochamiz
				if rotated(file, cfg.NginxLog, offset) {
					if reopened, openErr := os.Open(cfg.NginxLog); openErr == nil {
						file.Close()
						file, reader, offset = reopened, bufio.NewReader(reopened), 0
						slog.Info("nginx log was rotated, reopened", "path", cfg.NginxLog)
					}
				}

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

// rotated reports whether the file at path is no longer the one that is open (logrotate renamed it)
// or has become smaller than what was read (it was truncated).
func rotated(open *os.File, path string, offset int64) bool {
	current, err := os.Stat(path)
	if err != nil {
		// nginx has not created the new file yet: check again on the next pass.
		return false
	}

	opened, err := open.Stat()
	if err != nil {
		return true
	}

	return !os.SameFile(opened, current) || current.Size() < offset
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
		IPAddress: groups["ip"],
		Method:    groups["method"],
		Path:      groups["path"],
		Status:    status,
		Bytes:     bodyBytes,
		UserAgent: groups["ua"],
		LogTime:   logTime,
	}, nil
}
