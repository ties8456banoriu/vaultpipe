package secretrotate_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/nicholasgasior/vaultpipe/internal/audit"
	"github.com/nicholasgasior/vaultpipe/internal/secretrotate"
)

func TestRotate_AuditOnRecord(t *testing.T) {
	var buf bytes.Buffer
	logger, err := audit.NewLogger(&buf)
	if err != nil {
		t.Fatalf("audit.NewLogger: %v", err)
	}

	r := secretrotate.New()
	if err := r.Record("API_KEY", "secret/api", 3, secretrotate.PolicyScheduled); err != nil {
		t.Fatalf("Record: %v", err)
	}

	rec, err := r.Latest("API_KEY")
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}

	logger.Log(audit.Event{
		Action:  "secret_rotated",
		EnvKey:  rec.EnvKey,
		Details: map[string]any{"version": rec.Version, "policy": string(rec.Policy)},
	})

	var entry map[string]any
	if err := json.NewDecoder(&buf).Decode(&entry); err != nil {
		t.Fatalf("decode audit log: %v", err)
	}
	if entry["action"] != "secret_rotated" {
		t.Errorf("expected action secret_rotated, got %v", entry["action"])
	}
	if entry["env_key"] != "API_KEY" {
		t.Errorf("expected env_key API_KEY, got %v", entry["env_key"])
	}
}
