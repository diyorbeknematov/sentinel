package collector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/diyorbek/sentinel/agent/internal/config"
	"github.com/diyorbek/sentinel/agent/internal/models"
)

func nginxLine(path string) string {
	return fmt.Sprintf("1.2.3.4 - - [02/Jan/2026:15:04:05 +0000] \"GET %s HTTP/1.1\" 200 12 \"-\" \"curl\"\n", path)
}

func appendLine(t *testing.T, path, line string) {
	t.Helper()

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(line); err != nil {
		t.Fatal(err)
	}
}

func expectPath(t *testing.T, events <-chan models.Event, want string) {
	t.Helper()

	select {
	case e := <-events:
		got := e.Payload.(*models.NginxLogPayload).Path
		if got != want {
			t.Fatalf("got path %q, want %q", got, want)
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("no event for %q", want)
	}
}

func TestNginxCollectorSurvivesRotation(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "access.log")

	if err := os.WriteFile(logPath, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := make(chan models.Event, 10)
	go StartNginxLogCollector(ctx, &config.Config{NginxLog: logPath}, events)

	time.Sleep(300 * time.Millisecond) // the collector goes to the end of the file

	appendLine(t, logPath, nginxLine("/before"))
	expectPath(t, events, "/before")

	// logrotate: the file is renamed and nginx opens a new one.
	if err := os.Rename(logPath, logPath+".1"); err != nil {
		t.Fatal(err)
	}

	appendLine(t, logPath, nginxLine("/after-rename"))
	expectPath(t, events, "/after-rename")

	// copytruncate: the same file is emptied and written again. A file that is truncated and then
	// grows past the old offset before the next check cannot be told from a normal one, so the
	// collector has read a long line first and the new content is shorter than what it has read.
	long := "/" + strings.Repeat("a", 300)
	appendLine(t, logPath, nginxLine(long))
	expectPath(t, events, long)

	if err := os.Truncate(logPath, 0); err != nil {
		t.Fatal(err)
	}

	appendLine(t, logPath, nginxLine("/after-truncate"))
	expectPath(t, events, "/after-truncate")
}

func TestParseNginxLog(t *testing.T) {
	tests := []struct {
		name string
		line string
		want models.NginxLogPayload
		at   time.Time
	}{
		{
			name: "plain GET",
			line: `203.0.113.7 - - [02/Jan/2026:15:04:05 +0000] "GET /index.html HTTP/1.1" 200 612 "-" "Mozilla/5.0 (X11; Linux x86_64)"` + "\n",
			want: models.NginxLogPayload{
				IPAddress: "203.0.113.7", Method: "GET", Path: "/index.html", Status: 200, Bytes: 612,
				UserAgent: "Mozilla/5.0 (X11; Linux x86_64)",
			},
			at: time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC),
		},
		{
			name: "POST with a query, a referrer with a space and another time zone",
			line: `198.51.100.2 - - [31/Oct/2026:23:59:59 +0500] "POST /api/login?next=/home HTTP/2.0" 401 45 "https://example.com/a b" "curl/8.0"`,
			want: models.NginxLogPayload{
				IPAddress: "198.51.100.2", Method: "POST", Path: "/api/login?next=/home", Status: 401, Bytes: 45,
				UserAgent: "curl/8.0",
			},
			at: time.Date(2026, 10, 31, 23, 59, 59, 0, time.FixedZone("", 5*3600)),
		},
		{
			name: "IPv6 address, 404 and an empty user agent",
			line: `2001:db8::1 - - [02/Jan/2026:15:04:05 +0000] "GET /.env HTTP/1.1" 404 0 "-" "-"`,
			want: models.NginxLogPayload{
				IPAddress: "2001:db8::1", Method: "GET", Path: "/.env", Status: 404, Bytes: 0, UserAgent: "-",
			},
			at: time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNginxLog(tt.line)
			if err != nil {
				t.Fatal(err)
			}

			if !got.LogTime.Equal(tt.at) {
				t.Errorf("got time %v, want %v", got.LogTime, tt.at)
			}

			got.LogTime = time.Time{}

			if *got != tt.want {
				t.Errorf("got %+v, want %+v", *got, tt.want)
			}
		})
	}
}

func TestParseNginxLogRejectsBadLines(t *testing.T) {
	tests := map[string]string{
		"empty line":      "\n",
		"not a log":       "this is not an nginx line\n",
		"cut short":       `203.0.113.7 - - [02/Jan/2026:15:04:05 +0000] "GET /index.html`,
		"impossible date": `203.0.113.7 - - [99/Foo/2026:15:04:05 +0000] "GET / HTTP/1.1" 200 1 "-" "-"`,
	}

	for name, line := range tests {
		t.Run(name, func(t *testing.T) {
			if got, err := parseNginxLog(line); err == nil {
				t.Errorf("want an error, got %+v", *got)
			}
		})
	}
}

// startCollector runs the nginx collector on the file and returns its events and a function that
// stops it and waits until it has returned.
func startCollector(t *testing.T, path string) (<-chan models.Event, func()) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	events := make(chan models.Event, 10)
	done := make(chan struct{})

	go func() {
		StartNginxLogCollector(ctx, &config.Config{AgentID: "agent-1", NginxLog: path}, events)
		close(done)
	}()

	stop := func() {
		cancel()

		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Error("the collector did not stop after the context was cancelled")
		}
	}

	t.Cleanup(stop)

	return events, stop
}

func TestNginxCollectorReadsOnlyNewLines(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "access.log")

	// A line that was written before the agent started must not be sent again.
	if err := os.WriteFile(logPath, []byte(nginxLine("/old")), 0o644); err != nil {
		t.Fatal(err)
	}

	events, _ := startCollector(t, logPath)

	time.Sleep(300 * time.Millisecond)

	appendLine(t, logPath, nginxLine("/new"))

	select {
	case e := <-events:
		if e.Type != models.EventNginxLog || e.AgentID != "agent-1" {
			t.Errorf("got type %q and agent %q", e.Type, e.AgentID)
		}

		if got := e.Payload.(*models.NginxLogPayload).Path; got != "/new" {
			t.Errorf("got path %q, want /new (the old line must be skipped)", got)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("no event for the new line")
	}
}

func TestNginxCollectorSkipsBadLines(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "access.log")

	if err := os.WriteFile(logPath, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	events, _ := startCollector(t, logPath)

	time.Sleep(300 * time.Millisecond)

	appendLine(t, logPath, "this is not an nginx line\n")
	appendLine(t, logPath, nginxLine("/good"))
	appendLine(t, logPath, nginxLine("/also-good"))

	expectPath(t, events, "/good")
	expectPath(t, events, "/also-good")
}

func TestNginxCollectorStopsWithTheContext(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "access.log")

	if err := os.WriteFile(logPath, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	_, stop := startCollector(t, logPath)

	time.Sleep(300 * time.Millisecond)

	stop()
}

func TestNginxCollectorReturnsWhenTheFileIsMissing(t *testing.T) {
	events := make(chan models.Event, 1)
	done := make(chan struct{})

	go func() {
		StartNginxLogCollector(context.Background(), &config.Config{NginxLog: filepath.Join(t.TempDir(), "missing.log")}, events)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("the collector keeps running without a log file")
	}
}
