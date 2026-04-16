package telemetry_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/telemetry"
)

func TestStart_RecordsSpan(t *testing.T) {
	c := telemetry.New()
	s := c.Start("vault.fetch")
	if s.Name != "vault.fetch" {
		t.Fatalf("expected name vault.fetch, got %s", s.Name)
	}
	if s.Start.IsZero() {
		t.Fatal("expected non-zero start time")
	}
}

func TestStop_SetsEndAndDuration(t *testing.T) {
	c := telemetry.New()
	s := c.Start("env.write")
	time.Sleep(2 * time.Millisecond)
	c.Stop(s)

	if !s.Finished() {
		t.Fatal("expected span to be finished")
	}
	if s.Duration <= 0 {
		t.Fatalf("expected positive duration, got %v", s.Duration)
	}
	if s.Err != nil {
		t.Fatalf("expected nil error, got %v", s.Err)
	}
}

func TestStopWithError_RecordsError(t *testing.T) {
	c := telemetry.New()
	s := c.Start("vault.fetch")
	expErr := errors.New("connection refused")
	c.StopWithError(s, expErr)

	if s.Err == nil {
		t.Fatal("expected error to be recorded")
	}
	if !errors.Is(s.Err, expErr) {
		t.Fatalf("expected %v, got %v", expErr, s.Err)
	}
}

func TestSpans_ReturnsCopy(t *testing.T) {
	c := telemetry.New()
	c.Stop(c.Start("a"))
	c.Stop(c.Start("b"))

	spans := c.Spans()
	if len(spans) != 2 {
		t.Fatalf("expected 2 spans, got %d", len(spans))
	}
	// Mutating the returned slice should not affect the collector.
	spans[0].Name = "mutated"
	original := c.Spans()
	if original[0].Name == "mutated" {
		t.Fatal("Spans() should return a copy, not a reference")
	}
}

func TestWriteSummary_ContainsSpanNames(t *testing.T) {
	c := telemetry.New()
	c.Stop(c.Start("vault.fetch"))
	c.StopWithError(c.Start("env.write"), errors.New("disk full"))

	var buf bytes.Buffer
	c.WriteSummary(&buf)
	out := buf.String()

	if !strings.Contains(out, "vault.fetch") {
		t.Errorf("expected vault.fetch in summary, got:\n%s", out)
	}
	if !strings.Contains(out, "env.write") {
		t.Errorf("expected env.write in summary, got:\n%s", out)
	}
	if !strings.Contains(out, "disk full") {
		t.Errorf("expected error message in summary, got:\n%s", out)
	}
	if !strings.Contains(out, "ok") {
		t.Errorf("expected 'ok' status in summary, got:\n%s", out)
	}
}
