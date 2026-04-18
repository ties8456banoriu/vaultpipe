package secretsanitize_test

import (
	"testing"

	"github.com/vaultpipe/vaultpipe/internal/secretsanitize"
)

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	s := secretsanitize.New(secretsanitize.Rule{})
	_, err := s.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_TrimSpace(t *testing.T) {
	s := secretsanitize.New(secretsanitize.Rule{TrimSpace: true})
	out, err := s.Apply(map[string]string{"KEY": "  hello  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["KEY"]; got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestApply_StripPrefix(t *testing.T) {
	s := secretsanitize.New(secretsanitize.Rule{StripPrefix: "vault:"})
	out, err := s.Apply(map[string]string{"TOKEN": "vault:abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["TOKEN"]; got != "abc123" {
		t.Errorf("expected 'abc123', got %q", got)
	}
}

func TestApply_StripSuffix(t *testing.T) {
	s := secretsanitize.New(secretsanitize.Rule{StripSuffix: "=="})
	out, err := s.Apply(map[string]string{"KEY": "dGVzdA=="})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["KEY"]; got != "dGVzdA" {
		t.Errorf("expected 'dGVzdA', got %q", got)
	}
}

func TestApply_ReplaceNewlines(t *testing.T) {
	s := secretsanitize.New(secretsanitize.Rule{ReplaceNewlines: true})
	out, err := s.Apply(map[string]string{"CERT": "line1\nline2\r\nline3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["CERT"]; got != "line1 line2 line3" {
		t.Errorf("unexpected value: %q", got)
	}
}

func TestApply_NoRules_ReturnsUnchanged(t *testing.T) {
	s := secretsanitize.New(secretsanitize.Rule{})
	input := map[string]string{"A": "  val  ", "B": "other"}
	out, err := s.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range input {
		if out[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, out[k])
		}
	}
}

func TestApply_ReturnsCopy(t *testing.T) {
	s := secretsanitize.New(secretsanitize.Rule{TrimSpace: true})
	input := map[string]string{"X": " y "}
	out, _ := s.Apply(input)
	out["X"] = "mutated"
	if input["X"] != " y " {
		t.Error("Apply should not mutate input map")
	}
}
