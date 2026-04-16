package hooks_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/yourusername/vaultpipe/internal/audit"
	"github.com/yourusername/vaultpipe/internal/hooks"
)

func TestHooks_AuditOnFailure(t *testing.T) {
	var buf bytes.Buffer
	logger, err := audit.NewLogger(&buf)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	h := []hooks.Hook{{Stage: hooks.StagePre, Command: "false"}}
	r := hooks.NewRunner(h)

	hookErr := r.Run(context.Background(), hooks.StagePre)
	if hookErr == nil {
		t.Fatal("expected hook error")
	}

	logger.Log(audit.Entry{
		Event:   "hook_failed",
		Message: hookErr.Error(),
	})

	var entry map[string]interface{}
	if err := json.NewDecoder(&buf).Decode(&entry); err != nil {
		t.Fatalf("failed to decode audit log: %v", err)
	}
	if entry["event"] != "hook_failed" {
		t.Errorf("expected event hook_failed, got %v", entry["event"])
	}
	if entry["message"] == "" {
		t.Error("expected non-empty message")
	}
}
