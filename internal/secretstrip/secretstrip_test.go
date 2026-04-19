package secretstrip_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretstrip"
)

var base = map[string]string{
	"DB_PASSWORD":  "secret",
	"DB_HOST":      "localhost",
	"API_KEY":      "key123",
	"API_SECRET":   "topsecret",
	"APP_DEBUG":    "true",
}

func TestNew_NoPatterns_ReturnsError(t *testing.T) {
	_, err := secretstrip.New(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNew_EmptyPattern_ReturnsError(t *testing.T) {
	_, err := secretstrip.New([]string{""})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNew_InvalidGlob_ReturnsError(t *testing.T) {
	_, err := secretstrip.New([]string{"[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid glob, got nil")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	s, _ := secretstrip.New([]string{"DB_*"})
	_, err := s.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_RemovesMatchingKeys(t *testing.T) {
	s, err := secretstrip.New([]string{"DB_*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	, err := s.Apply(base)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("expected DB_PASSWORD to be stripped")
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be stripped")
	}
	if _, ok := out["API_KEY"]; !ok {
		t.Error("expected API_KEY to be retained")
	}
}

func TestApply_MultiplePatterns(t *testing.T) {
	s, err := secretstrip.New([]string{"DB_*", "API_*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := s.Apply(base)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 key remaining, got %d", len(out))
	}
	if _, ok := out["APP_DEBUG"]; !ok {
		t.Error("expected APP_DEBUG to be retained")
	}
}

func TestApply_NoMatchRetainsAll(t *testing.T) {
	s, err := secretstrip.New([]string{"NOOP_*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := s.Apply(base)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if len(out) != len(base) {
		t.Fatalf("expected all %d keys, got %d", len(base), len(out))
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	s, _ := secretstrip.New([]string{"DB_*"})
	input := map[string]string{"DB_HOST": "localhost", "APP_NAME": "vaultpipe"}
	s.Apply(input) //nolint
	if _, ok := input["DB_HOST"]; !ok {
		t.Error("Apply must not mutate the input map")
	}
}
