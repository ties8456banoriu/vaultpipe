package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/audit"
)

func newBufLogger() (*audit.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	return audit.NewLogger(&buf), &buf
}

func TestLog_WritesJSONLine(t *testing.T) {
	l, buf := newBufLogger()

	err := l.Log(audit.Event{
		Type:       audit.EventSecretFetched,
		SecretPath: "secret/data/myapp",
		KeyCount:   3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var got audit.Event
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got.Type != audit.EventSecretFetched {
		t.Errorf("expected type %q, got %q", audit.EventSecretFetched, got.Type)
	}
	if got.SecretPath != "secret/data/myapp" {
		t.Errorf("expected path %q, got %q", "secret/data/myapp", got.SecretPath)
	}
	if got.KeyCount != 3 {
		t.Errorf("expected key_count 3, got %d", got.KeyCount)
	}
}

func TestLog_SetsTimestampIfZero(t *testing.T) {
	l, buf := newBufLogger()
	before := time.Now().UTC()
	_ = l.Log(audit.Event{Type: audit.EventEnvWritten})
	after := time.Now().UTC()

	var got audit.Event
	_ = json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &got)

	if got.Timestamp.Before(before) || got.Timestamp.After(after) {
		t.Errorf("timestamp %v not within expected range [%v, %v]", got.Timestamp, before, after)
	}
}

func TestLogSecretFetched(t *testing.T) {
	l, buf := newBufLogger()
	_ = l.LogSecretFetched("secret/data/db", 5)

	var got audit.Event
	_ = json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &got)

	if got.Type != audit.EventSecretFetched {
		t.Errorf("expected %q, got %q", audit.EventSecretFetched, got.Type)
	}
	if got.KeyCount != 5 {
		t.Errorf("expected key_count 5, got %d", got.KeyCount)
	}
}

func TestLogEnvWritten(t *testing.T) {
	l, buf := newBufLogger()
	_ = l.LogEnvWritten(".env", 4)

	var got audit.Event
	_ = json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &got)

	if got.Type != audit.EventEnvWritten {
		t.Errorf("expected %q, got %q", audit.EventEnvWritten, got.Type)
	}
	if got.OutputFile != ".env" {
		t.Errorf("expected output_file '.env', got %q", got.OutputFile)
	}
}

func TestLogError(t *testing.T) {
	l, buf := newBufLogger()
	_ = l.LogError("vault unreachable")

	var got audit.Event
	_ = json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &got)

	if got.Type != audit.EventError {
		t.Errorf("expected %q, got %q", audit.EventError, got.Type)
	}
	if got.Message != "vault unreachable" {
		t.Errorf("expected message 'vault unreachable', got %q", got.Message)
	}
}
