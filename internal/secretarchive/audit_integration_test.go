package secretarchive_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/audit"
	"github.com/yourusername/vaultpipe/internal/secretarchive"
)

func TestArchive_AuditOnStore(t *testing.T) {
	var buf bytes.Buffer
	logger, err := audit.NewLogger(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a := secretarchive.New()
	secrets := map[string]string{"TOKEN": "abc", "DB": "pass"}
	if err := a.Store("release-1", secrets); err != nil {
		t.Fatalf("store failed: %v", err)
	}

	e, err := a.Get("release-1")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	logger.Log(audit.Event{
		Action:    "archive.store",
		Timestamp: time.Now().UTC(),
		Details: map[string]string{
			"name": e.Name,
			"keys": fmt.Sprintf("%d", len(e.Secrets)),
		},
	})

	if buf.Len() == 0 {
		t.Fatal("expected audit log output")
	}
	var entry map[string]interface{}
	if err := json.NewDecoder(&buf).Decode(&entry); err != nil {
		t.Fatalf("invalid JSON in audit log: %v", err)
	}
	if entry["action"] != "archive.store" {
		t.Errorf("expected action archive.store, got %v", entry["action"])
	}
}
