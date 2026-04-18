package secretnamespace_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/audit"
	"github.com/yourusername/vaultpipe/internal/secretnamespace"
)

func TestNamespace_AuditOnMerge(t *testing.T) {
	var buf bytes.Buffer
	logger, err := audit.NewLogger(&buf)
	if err != nil {
		t.Fatalf("failed to create audit logger: %v", err)
	}

	m := secretnamespace.New()
	_ = m.Set("alpha", map[string]string{"TOKEN": "abc"})
	_ = m.Set("beta", map[string]string{"SECRET": "xyz"})

	result, err := m.Merge("alpha", "beta")
	if err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	logger.Log(audit.Event{
		Timestamp: time.Now(),
		Action:    "namespace_merge",
		Detail:    fmt.Sprintf("merged %d keys from namespaces: alpha, beta", len(result)),
	})

	var entry map[string]interface{}
	if err := json.NewDecoder(&buf).Decode(&entry); err != nil {
		t.Fatalf("failed to decode audit log: %v", err)
	}
	if entry["action"] != "namespace_merge" {
		t.Errorf("expected action namespace_merge, got %v", entry["action"])
	}
}
