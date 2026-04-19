package secretmask_test

import (
	"testing"

	"github.com/eliziario/vaultpipe/internal/secretmask"
)

func TestParseRules_Valid(t *testing.T) {
	rules, err := secretmask.ParseRules([]string{`\d+=***`, `secret=REDACTED`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
}

func TestParseRules_InvalidFormat(t *testing.T) {
	_, err := secretmask.ParseRules([]string{"noequals"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseRules_InvalidRegex(t *testing.T) {
	_, err := secretmask.ParseRules([]string{`[invalid=replacement`})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	m := secretmask.New(nil)
	_, err := m.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_NoRules_ReturnsUnchanged(t *testing.T) {
	m := secretmask.New(nil)
	in := map[string]string{"KEY": "value123"}
	out, err := m.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value123" {
		t.Errorf("expected 'value123', got '%s'", out["KEY"])
	}
}

func TestApply_MasksMatchingValues(t *testing.T) {
	rules, _ := secretmask.ParseRules([]string{`\d+=***`})
	m := secretmask.New(rules)
	in := map[string]string{"TOKEN": "abc123", "NAME": "alice"}
	out, err := m.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "abc***" {
		t.Errorf("expected 'abc***', got '%s'", out["TOKEN"])
	}
	if out["NAME"] != "alice" {
		t.Errorf("expected 'alice', got '%s'", out["NAME"])
	}
}

func TestApply_MultipleRulesAppliedInOrder(t *testing.T) {
	rules, _ := secretmask.ParseRules([]string{`secret=REDACTED`, `REDACTED=HIDDEN`})
	m := secretmask.New(rules)
	in := map[string]string{"KEY": "mysecretvalue"}
	out, err := m.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "myREDACTEDvalue" {
		t.Errorf("unexpected value: %s", out["KEY"])
	}
}

func TestApply_ReturnsCopy(t *testing.T) {
	m := secretmask.New(nil)
	in := map[string]string{"A": "1"}
	out, _ := m.Apply(in)
	out["A"] = "modified"
	if in["A"] != "1" {
		t.Error("Apply modified the original map")
	}
}
