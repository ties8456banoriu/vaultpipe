package validate_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/audit"
	"github.com/yourusername/vaultpipe/internal/validate"
)

func TestValidate_AuditWarningsLogged(t *testing.T) {
	var buf bytes.Buffer
	logger, err := audit.NewLogger(&buf)
	if err != nil {
		t.Fatalf("failed to create audit logger: %v", err)
	}

	v := validate.NewValidator(true)
	secrets := map[string]string{
		"PRESENT_KEY": "value",
		"EMPTY_KEY":   "",
	}

	warnings, valErr := v.Validate(secrets)
	if valErr != nil {
		t.Fatalf("unexpected validation error: %v", valErr)
	}

	for _, w := range warnings {
		logger.Log(audit.Entry{
			Timestamp: time.Now(),
			Event:     "secret_warning",
			Key:       w.Key,
			Detail:    w.Warning,
		})
	}

	var entry map[string]interface{}
	if err := json.NewDecoder(&buf).Decode(&entry); err != nil {
		t.Fatalf("failed to decode audit log entry: %v", err)
	}
	if entry["event"] != "secret_warning" {
		t.Errorf("expected event=secret_warning, got %v", entry["event"])
	}
	if entry["key"] != "EMPTY_KEY" {
		t.Errorf("expected key=EMPTY_KEY, got %v", entry["key"])
	}
}
