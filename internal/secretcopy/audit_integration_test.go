package secretcopy_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/audit"
	"github.com/yourusername/vaultpipe/internal/secretcopy"
)

func TestCopy_AuditAfterCopy(t *testing.T) {
	var buf bytes.Buffer
	logger, err := audit.NewLogger(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c, err := secretcopy.New([]secretcopy.Rule{{From: "DB_HOST", To: "PG_HOST"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets := map[string]string{"DB_HOST": "localhost"}
	out, err := c.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for k := range out {
		logger.LogSecretFetched(k, "secret/data/db", time.Now())
	}

	var entry map[string]interface{}
	if err := json.NewDecoder(&buf).Decode(&entry); err != nil {
		t.Fatalf("failed to decode audit log: %v", err)
	}
	if entry["event"] != "secret_fetched" {
		t.Errorf("expected event=secret_fetched, got %v", entry["event"])
	}
}
