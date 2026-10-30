package collector

import (
	"context"
	"testing"
	"time"

	"github.com/diyorbek/sentinel/agent/internal/config"
	"github.com/diyorbek/sentinel/agent/internal/models"
)

func TestCollectMetrics(t *testing.T) {
	metric, err := collectMetrics()
	if err != nil {
		t.Fatal(err)
	}

	for name, value := range map[string]float64{"cpu": metric.CPU, "ram": metric.RAM, "disk": metric.Disk} {
		if value < 0 || value > 100 {
			t.Errorf("%s is %.1f, want a percent between 0 and 100", name, value)
		}
	}

	if time.Since(metric.LogTime) > 5*time.Second {
		t.Errorf("log time %v is not recent", metric.LogTime)
	}
}

func TestMetricsCollectorSendsEventsAndStops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := &config.Config{AgentID: "agent-1", MetricsInterval: 50 * time.Millisecond}
	events := make(chan models.Event, 1)
	done := make(chan struct{})

	go func() {
		StartMetricsCollector(ctx, cfg, events)
		close(done)
	}()

	select {
	case e := <-events:
		if e.Type != models.EventMetric {
			t.Errorf("got type %q, want %q", e.Type, models.EventMetric)
		}

		if e.AgentID != "agent-1" {
			t.Errorf("got agent %q, want agent-1", e.AgentID)
		}

		if _, ok := e.Payload.(*models.MetricPayload); !ok {
			t.Errorf("payload is %T, want *models.MetricPayload", e.Payload)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("no metric event")
	}

	cancel()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("the collector did not stop after the context was cancelled")
	}
}

func TestMetricsCollectorSendsEventsRepeatedly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &config.Config{AgentID: "agent-1", MetricsInterval: 20 * time.Millisecond}
	events := make(chan models.Event, 5)

	go StartMetricsCollector(ctx, cfg, events)

	var last time.Time

	for i := 0; i < 3; i++ {
		select {
		case e := <-events:
			at := e.Payload.(*models.MetricPayload).LogTime
			if !at.After(last) {
				t.Errorf("event %d has time %v, not after the previous %v", i, at, last)
			}

			last = at
		case <-time.After(5 * time.Second):
			t.Fatalf("got %d events, want 3", i)
		}
	}
}

func TestMetricsCollectorStopsBeforeTheFirstTick(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})

	go func() {
		StartMetricsCollector(ctx, &config.Config{MetricsInterval: time.Hour}, make(chan models.Event))
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("the collector did not stop")
	}
}

// With nobody reading the channel the collector blocks on sending; cancelling the context must free it.
func TestMetricsCollectorStopsWhileBlockedOnSend(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := &config.Config{MetricsInterval: 10 * time.Millisecond}
	done := make(chan struct{})

	go func() {
		StartMetricsCollector(ctx, cfg, make(chan models.Event)) // no reader
		close(done)
	}()

	// Long enough for the first metric (about a second of CPU sampling) to be ready and stuck.
	time.Sleep(1500 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("the collector is stuck on sending after the context was cancelled")
	}
}
