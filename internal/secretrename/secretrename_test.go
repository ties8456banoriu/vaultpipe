package secretrename_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/secretrename"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "abc123",
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := secretrename.New(nil)
	if err == nil {
		t.Fatal("expected error for no rules")
	}
}

func TestNew_EmptyFrom_ReturnsError(t *testing.T) {
	_, err := secretrename.New([]secretrename.Rule{{From: "", To: "NEW_KEY"}})
	if err == nil {
		t.Fatal("expected error for empty From")
	}
}

func TestNew_EmptyTo_ReturnsError(t *testing.T) {
	_, err := secretrename.New([]secretrename.Rule{{From: "DB_HOST", To: ""}})
	if err == nil {
		t.Fatal("expected error for empty To")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	r, _ := secretrename.New([]secretrename.Rule{{From: "A", To: "B"}})
	_, err := r.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_RenamesKey(t *testing.T) {
	r, err := secretrename.New([]secretrename.Rule{{From: "DB_HOST", To: "DATABASE_HOST"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := r.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be removed")
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
}

func TestApply_UnknownKey_ReturnsError(t *testing.T) {
	r, _ := secretrename.New([]secretrename.Rule{{From: "MISSING_KEY", To: "NEW_KEY"}})
	_, err := r.Apply(baseSecrets())
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	rules := []secretrename.Rule{
		{From: "DB_HOST", To: "DATABASE_HOST"},
		{From: "DB_PORT", To: "DATABASE_PORT"},
	}
	r, _ := secretrename.New(rules)
	out, err := r.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATABASE_HOST"] != "localhost" || out["DATABASE_PORT"] != "5432" {
		t.Error("expected both keys to be renamed")
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("old key DB_HOST should not exist")
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := secretrename.ParseRules([]string{"DB_HOST:DATABASE_HOST", "API_KEY:SERVICE_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].From != "DB_HOST" || rules[0].To != "DATABASE_HOST" {
		t.Errorf("unexpected rule: %+v", rules[0])
	}
}

func TestParseRules_Invalid(t *testing.T) {
	_, err := secretrename.ParseRules([]string{"INVALID_NO_COLON"})
	if err == nil {
		t.Fatal("expected error for invalid rule format")
	}
}
