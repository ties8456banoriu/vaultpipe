package mapping_test

import (
	"testing"

	"github.com/yourorg/vaultpipe/internal/mapping"
)

func TestApply_DefaultKeys(t *testing.T) {
	mapper := mapping.NewMapper(nil)
	secrets := map[string]string{
		"db-password": "s3cr3t",
		"api.key":     "abc123",
	}

	result, err := mapper.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("expected DB_PASSWORD=s3cr3t, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", result["API_KEY"])
	}
}

func TestApply_WithRules(t *testing.T) {
	rules := []mapping.Rule{
		{VaultKey: "db-password", EnvKey: "DATABASE_PASSWORD"},
	}
	mapper := mapping.NewMapper(rules)
	secrets := map[string]string{
		"db-password": "hunter2",
		"api-token":   "tok123",
	}

	result, err := mapper.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["DATABASE_PASSWORD"] != "hunter2" {
		t.Errorf("expected DATABASE_PASSWORD=hunter2, got %q", result["DATABASE_PASSWORD"])
	}
	// api-token has no rule, should fall back to default
	if result["API_TOKEN"] != "tok123" {
		t.Errorf("expected API_TOKEN=tok123, got %q", result["API_TOKEN"])
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	mapper := mapping.NewMapper(nil)
	_, err := mapper.Apply(map[string]string{})
	if err == nil {
		t.Error("expected error for empty secrets, got nil")
	}
}

func TestParseRules_Valid(t *testing.T) {
	raw := []string{"db-password=DB_PASS", "api.key=API_KEY"}
	rules, err := mapping.ParseRules(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].VaultKey != "db-password" || rules[0].EnvKey != "DB_PASS" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestParseRules_Invalid(t *testing.T) {
	cases := []string{"nodivider", "=nokey", "noval="}
	for _, c := range cases {
		_, err := mapping.ParseRules([]string{c})
		if err == nil {
			t.Errorf("expected error for rule %q, got nil", c)
		}
	}
}
