package secretomit_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretomit"
)

func TestNew_EmptyPatterns_ReturnsError(t *testing.T) {
	_, err := secretomit.New([]string{})
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestNew_EmptyPatternString_ReturnsError(t *testing.T) {
	_, err := secretomit.New([]string{""})
	if err == nil {
		t.Fatal("expected error for empty pattern string")
	}
}

func TestNew_InvalidRegex_ReturnsError(t *testing.T) {
	_, err := secretomit.New([]string{"[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	o, err := secretomit.New([]string{"^$"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = o.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_RemovesEmptyValues(t *testing.T) {
	o, _ := secretomit.New([]string{"^$"})
	secrets := map[string]string{"A": "hello", "B": "", "C": "world"}
	out, err := o.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["B"]; ok {
		t.Error("expected B to be omitted")
	}
	if out["A"] != "hello" || out["C"] != "world" {
		t.Error("expected A and C to be retained")
	}
}

func TestApply_RemovesPlaceholders(t *testing.T) {
	o, _ := secretomit.New([]string{"^CHANGEME$", "^TODO$"})
	secrets := map[string]string{"X": "CHANGEME", "Y": "real-value", "Z": "TODO"}
	out, err := o.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["X"]; ok {
		t.Error("expected X to be omitted")
	}
	if _, ok := out["Z"]; ok {
		t.Error("expected Z to be omitted")
	}
	if out["Y"] != "real-value" {
		t.Error("expected Y to be retained")
	}
}

func TestApply_NoMatchRetainsAll(t *testing.T) {
	o, _ := secretomit.New([]string{"^PLACEHOLDER$"})
	secrets := map[string]string{"A": "foo", "B": "bar"}
	out, err := o.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}
