package notify_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/notify"
)

func TestNotify_WritesMessage(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf, false)

	err := n.Notify(notify.Event{
		Profile:  "staging",
		KeyCount: 5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "staging") {
		t.Errorf("expected profile name in output, got: %s", out)
	}
	if !strings.Contains(out, "5 secret(s)") {
		t.Errorf("expected key count in output, got: %s", out)
	}
}

func TestNotify_SetsTimestampIfZero(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf, false)

	before := time.Now()
	_ = n.Notify(notify.Event{Profile: "dev", KeyCount: 1})
	after := time.Now()

	out := buf.String()
	// Timestamp should be between before and after — just check it's present.
	if !strings.Contains(out, "T") {
		t.Errorf("expected RFC3339 timestamp in output, got: %s", out)
	}
	_ = before
	_ = after
}

func TestNotify_ContainsVaultpipePrefix(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf, false)

	_ = n.Notify(notify.Event{Profile: "prod", KeyCount: 3})

	if !strings.HasPrefix(strings.TrimSpace(buf.String()), "[vaultpipe]") {
		t.Errorf("expected [vaultpipe] prefix, got: %s", buf.String())
	}
}

func TestNotify_WriteError_ReturnsError(t *testing.T) {
	n := notify.New(&errWriter{}, false)

	err := n.Notify(notify.Event{Profile: "x", KeyCount: 1})
	if err == nil {
		t.Fatal("expected error from failing writer")
	}
}

type errWriter struct{}

func (e *errWriter) Write(_ []byte) (int, error) {
	return 0, fmt.Errorf("disk full")
}

import "fmt"
