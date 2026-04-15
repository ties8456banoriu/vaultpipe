package diff_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/yourorg/vaultpipe/internal/audit"
	"github.com/yourorg/vaultpipe/internal/diff"
)

// TestDiffAuditIntegration verifies that diff results can be serialised
// and emitted through the audit logger without data loss.
func TestDiffAuditIntegration(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	prev := map[string]string{"DB_PASS": "old", "API_KEY": "same"}
	curr := map[string]string{"API_KEY": "same", "NEW_SECRET": "val"}

	res := diff.Compare(prev, curr)

	logger.Log(audit.Event{
		Timestamp: time.Now(),
		Action:    "secrets_diff",
		Meta: map[string]interface{}{
			"added":   res.Added,
			"removed": res.Removed,
			"changed": res.Changed,
		},
	})

	var entry map[string]interface{}
	if err := json.NewDecoder(&buf).Decode(&entry); err != nil {
		t.Fatalf("failed to decode audit log entry: %v", err)
	}

	if entry["action"] != "secrets_diff" {
		t.Errorf("unexpected action: %v", entry["action"])
	}

	meta, ok := entry["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("meta field missing or wrong type")
	}

	added, _ := meta["added"].([]interface{})
	if len(added) != 1 || added[0] != "NEW_SECRET" {
		t.Errorf("expected added=[NEW_SECRET] in audit meta, got %v", added)
	}

	removed, _ := meta["removed"].([]interface{})
	if len(removed) != 1 || removed[0] != "DB_PASS" {
		t.Errorf("expected removed=[DB_PASS] in audit meta, got %v", removed)
	}
}
