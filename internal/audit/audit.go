// Package audit provides lightweight audit logging for secret injection events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventType represents the kind of audit event.
type EventType string

const (
	EventSecretFetched  EventType = "secret_fetched"
	EventEnvWritten     EventType = "env_written"
	EventError          EventType = "error"
)

// Event holds metadata about a single audit log entry.
type Event struct {
	Timestamp  time.Time `json:"timestamp"`
	Type       EventType `json:"type"`
	SecretPath string    `json:"secret_path,omitempty"`
	OutputFile string    `json:"output_file,omitempty"`
	KeyCount   int       `json:"key_count,omitempty"`
	Message    string    `json:"message,omitempty"`
}

// Logger writes structured audit events as JSON lines.
type Logger struct {
	out io.Writer
}

// NewLogger creates a Logger that writes to the given writer.
// Pass nil to use stderr as the default destination.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{out: w}
}

// Log encodes and writes a single Event as a JSON line.
func (l *Logger) Log(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintln(l.out, string(data))
	return err
}

// LogSecretFetched is a convenience helper for secret fetch events.
func (l *Logger) LogSecretFetched(path string, keyCount int) error {
	return l.Log(Event{
		Type:       EventSecretFetched,
		SecretPath: path,
		KeyCount:   keyCount,
	})
}

// LogEnvWritten is a convenience helper for env-file write events.
func (l *Logger) LogEnvWritten(outputFile string, keyCount int) error {
	return l.Log(Event{
		Type:       EventEnvWritten,
		OutputFile: outputFile,
		KeyCount:   keyCount,
	})
}

// LogError is a convenience helper for error events.
func (l *Logger) LogError(msg string) error {
	return l.Log(Event{
		Type:    EventError,
		Message: msg,
	})
}
