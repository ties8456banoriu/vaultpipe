package secretsplit_test

import (
	"testing"

	"github.com/fmjstudios/vaultpipe/internal/secretsplit"
)

var baseSecrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASS":     "secret",
	"APP_KEY":     "abc123",
	"APP_SECRET":  "xyz789",
	"OTHER_VALUE": "misc",
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	s := secretsplit.New()
	_, err := s.Apply(map[string]string{}, []secretsplit.Rule{{Name: "db", Prefix: "DB_"}})
	if err != secretsplit.ErrEmptySecrets {
		t.Fatalf("expected ErrEmptySecrets, got %v", err)
	}
}

func TestApply_NoRules_ReturnsError(t *testing.T) {
	s := secretsplit.New()
	_, err := s.Apply(baseSecrets, nil)
	if err != secretsplit.ErrNoRules {
		t.Fatalf("expected ErrNoRules, got %v", err)
	}
}

func TestApply_EmptyRuleName_ReturnsError(t *testing.T) {
	s := secretsplit.New()
	_, err := s.Apply(baseSecrets, []secretsplit.Rule{{Name: "", Prefix: "DB_"}})
	if err != secretsplit.ErrEmptyName {
		t.Fatalf("expected ErrEmptyName, got %v", err)
	}
}

func TestApply_SplitsByPrefix(t *testing.T) {
	s := secretsplit.New()
	rules := []secretsplit.Rule{
		{Name: "db", Prefix: "DB_"},
		{Name: "app", Prefix: "APP_"},
	}
	result, err := s.Apply(baseSecrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result["db"]) != 2 {
		t.Errorf("expected 2 db secrets, got %d", len(result["db"]))
	}
	if len(result["app"]) != 2 {
		t.Errorf("expected 2 app secrets, got %d", len(result["app"]))
	}
}

func TestApply_UnmatchedKeysGrouped(t *testing.T) {
	s := secretsplit.New()
	rules := []secretsplit.Rule{
		{Name: "db", Prefix: "DB_"},
	}
	result, err := s.Apply(baseSecrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	unmatched, ok := result["_unmatched"]
	if !ok {
		t.Fatal("expected _unmatched subset")
	}
	if _, has := unmatched["OTHER_VALUE"]; !has {
		t.Error("expected OTHER_VALUE in _unmatched")
	}
}

func TestApply_CaseInsensitivePrefix(t *testing.T) {
	s := secretsplit.New()
	rules := []secretsplit.Rule{
		{Name: "db", Prefix: "db_"},
	}
	result, err := s.Apply(baseSecrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result["db"]) != 2 {
		t.Errorf("expected 2 db secrets with case-insensitive match, got %d", len(result["db"]))
	}
}
