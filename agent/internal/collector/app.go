package collector

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/diyorbek/sentinel/agent/internal/config"
	"github.com/diyorbek/sentinel/agent/internal/models"
)

// text log prefixlari — skip qilinadi

var skipPrefixes = []string{
	"[GIN]",       // Gin HTTP framework
	"[GIN-debug]", // Gin debug
}

func shouldSkip(line string) bool {
	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

func StartAppLogCollector(ctx context.Context, cfg *config.Config, eventCh chan<- models.Event) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		slog.Error("docker client failed", "err", err)
		return
	}

	slog.Info("starting app log collector")

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		slog.Error("container list failed", "err", err)
		cli.Close()
		return
	}

	for _, c := range containers {
		name := strings.TrimPrefix(c.Names[0], "/")
		if !isMonitored(name, cfg.AppContainers) {
			continue
		}
		slog.Info("monitoring container", "name", name)
		go streamContainerLogs(ctx, cli, c.ID, name, cfg, eventCh)
	}

	go watchNewContainers(ctx, cli, cfg, eventCh)

	<-ctx.Done()
	cli.Close()
}

func streamContainerLogs(
	ctx context.Context,
	cli *client.Client,
	containerID, name string,
	cfg *config.Config,
	eventCh chan<- models.Event,
) {
	for {
		if ctx.Err() != nil {
			return
		}

		if err := readLogs(ctx, cli, containerID, name, cfg, eventCh); err != nil {
			slog.Error("stream stopped", "container", name, "err", err)
		}

		if ctx.Err() != nil {
			return
		}

		slog.Info("reconnecting", "container", name)
		select {
		case <-time.After(2 * time.Second):
		case <-ctx.Done():
			return
		}
	}
}

func readLogs(
	ctx context.Context,
	cli *client.Client,
	containerID, name string,
	cfg *config.Config,
	eventCh chan<- models.Event,
) error {
	reader, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "0",
	})
	if err != nil {
		return fmt.Errorf("container logs: %w", err)
	}
	defer reader.Close()

	pr, pw := io.Pipe()
	copyErr := make(chan error, 1)

	go func() {
		_, err := stdcopy.StdCopy(pw, pw, reader)
		pw.CloseWithError(err)
		copyErr <- err
	}()

	scanner := bufio.NewScanner(pr)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || shouldSkip(line) {
			continue
		}

		appLog, err := parseAppLog(line)
		if err != nil {
			slog.Debug("skip unparseable line", "container", name, "line", line)
			continue
		}

		appLog.ServiceName = name
		fmt.Printf("app log: %s\n", line)
		select {
		case eventCh <- models.Event{
			Type:    models.EventAppLog,
			AgentID: cfg.AgentID,
			Payload: appLog,
		}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}

	return <-copyErr
}

func watchNewContainers(
	ctx context.Context,
	cli *client.Client,
	cfg *config.Config,
	eventCh chan<- models.Event,
) {
	for {
		if ctx.Err() != nil {
			return
		}

		if err := listenDockerEvents(ctx, cli, cfg, eventCh); err != nil {
			slog.Error("docker events error", "err", err)
		}

		if ctx.Err() != nil {
			return
		}

		select {
		case <-time.After(2 * time.Second):
		case <-ctx.Done():
			return
		}
	}
}

func listenDockerEvents(
	ctx context.Context,
	cli *client.Client,
	cfg *config.Config,
	eventCh chan<- models.Event,
) error {
	f := filters.NewArgs()
	f.Add("type", "container")
	f.Add("event", "start")

	eventsCh, errs := cli.Events(ctx, events.ListOptions{Filters: f})

	for {
		select {
		case e := <-eventsCh:
			name := e.Actor.Attributes["name"]
			if !isMonitored(name, cfg.AppContainers) {
				continue
			}
			slog.Info("new container detected", "container", name)
			go streamContainerLogs(ctx, cli, e.Actor.ID, name, cfg, eventCh)

		case err := <-errs:
			return err

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func isMonitored(name string, list []string) bool {
	for _, m := range list {
		if strings.Contains(name, m) {
			return true
		}
	}
	return false
}

func parseAppLog(line string) (*models.AppLogPayload, error) {
	if isJSON(line) {
		return parseJSON(line)
	}
	return parsePlainText(line)
}

func isJSON(line string) bool {
	return json.Valid([]byte(line))
}

func parseJSON(line string) (*models.AppLogPayload, error) {
	var raw models.RawLog
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("json parse: %w", err)
	}

	log := &models.AppLogPayload{
		Metadata: make(map[string]any),
	}

	for k, v := range raw {
		switch strings.ToLower(k) {
		case "level", "log_level", "severity":
			if s, ok := v.(string); ok {
				log.Level = s
			}
		case "event":
			if s, ok := v.(string); ok {
				log.Event = s
			}
		case "message", "msg":
			if s, ok := v.(string); ok {
				log.Message = s
			}
		case "time", "timestamp", "log_time":
			if ts, ok := v.(string); ok {
				for _, format := range models.TimeFormats {
					if t, err := time.Parse(format, ts); err == nil {
						log.LogTime = t
						break
					}
				}
			}
		default:
			log.Metadata[k] = v
		}
		if log.Event == "" {
			log.Event = "unknown_event"
		}
	}

	return log, nil
}

func parsePlainText(line string) (*models.AppLogPayload, error) {
	// plain text loglarni raw holda saqlaymiz
	// level ni keyword orqali aniqlaymiz
	return &models.AppLogPayload{
		Message: line,
		Level:   detectLevel(line),
		Event:   "plain_text_log",
		LogTime: time.Now().UTC(),
	}, nil
}

func detectLevel(line string) string {
	u := strings.ToUpper(line)
	switch {
	case strings.Contains(u, "ERROR") || strings.Contains(u, "FATAL"):
		return "ERROR"
	case strings.Contains(u, "WARN") || strings.Contains(u, "WARNING"):
		return "WARN"
	case strings.Contains(u, "DEBUG"):
		return "DEBUG"
	default:
		return "INFO"
	}
}
