package secretclone_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/audit"
	"github.com/yourusername/vaultpipe/internal/secretclone"
)

func TestClone_AuditAfterClone(t *testing.T) {
	var buf bytes.Buffer
	logger, err := audit.NewLogger(&buf)
	if err != nil {
		t.Fatalf("failed to create audit logger: %v", err)
	}

	src := map[string]string{"api_key": "topsecret", "db_pass": "hunter2"}
	c := secretclone.New(secretclone.WithUppercase())
	out, err := c.Clone(src)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	for k := range out {
		logger.Log(audit.Entry{
			Timestamp: time.Now(),
			Event:     "secret_cloned",
			Detail:    k,
		})
	}

	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	if len(lines) != len(out) {
		t.Fatalf("expected %d audit lines, got %d", len(out), len(lines))
	}
	for _, line := range lines {
		var entry map[string]interface{}
		if err := json.Unmarshal(line, &entry); err != nil {
			t.Fatalf("invalid JSON audit line: %s", line)
		}
		if entry["event"] != "secret_cloned" {
			t.Errorf("unexpected event: %v", entry["event"])
		}
	}
}
