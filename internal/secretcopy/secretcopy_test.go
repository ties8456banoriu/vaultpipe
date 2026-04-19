package secretcopy_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretcopy"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret123",
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := secretcopy.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptyFrom_ReturnsError(t *testing.T) {
	_, err := secretcopy.New([]secretcopy.Rule{{From: "", To: "X"}})
	if err == nil {
		t.Fatal("expected error for empty From")
	}
}

func TestNew_EmptyTo_ReturnsError(t *testing.T) {
	_, err := secretcopy.New([]secretcopy.Rule{{From: "A", To: ""}})
	if err == nil {
		t.Fatal("expected error for empty To")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	c, _ := secretcopy.New([]secretcopy.Rule{{From: "A", To: "B"}})
	_, err := c.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_CopiesKey(t *testing.T) {
	c, err := secretcopy.New([]secretcopy.Rule{{From: "DB_HOST", To: "DATABASE_HOST"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := c.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
	if out["DB_HOST"] != "localhost" {
		t.Error("original key DB_HOST should be preserved")
	}
}

func TestApply_MissingSourceKey_ReturnsError(t *testing.T) {
	c, _ := secretcopy.New([]secretcopy.Rule{{From: "MISSING", To: "TARGET"}})
	_, err := c.Apply(baseSecrets())
	if err == nil {
		t.Fatal("expected error for missing source key")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	c, err := secretcopy.New([]secretcopy.Rule{
		{From: "DB_HOST", To: "PG_HOST"},
		{From: "DB_PORT", To: "PG_PORT"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := c.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PG_HOST"] != "localhost" {
		t.Errorf("expected PG_HOST=localhost, got %q", out["PG_HOST"])
	}
	if out["PG_PORT"] != "5432" {
		t.Errorf("expected PG_PORT=5432, got %q", out["PG_PORT"])
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := secretcopy.ParseRules([]string{"A:B", "C:D"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].From != "A" || rules[0].To != "B" {
		t.Errorf("unexpected rule: %+v", rules[0])
	}
}

func TestParseRules_Invalid(t *testing.T) {
	_, err := secretcopy.ParseRules([]string{"NODIVIDER"})
	if err == nil {
		t.Fatal("expected error for invalid rule format")
	}
}
