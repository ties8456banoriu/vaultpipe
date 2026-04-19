package secrettrim_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secrettrim"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"TOKEN": "Bearer abc123xyz",
		"SHORT": "hi",
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := secrettrim.New(nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNew_NegativeStart_ReturnsError(t *testing.T) {
	_, err := secrettrim.New([]secrettrim.Rule{{Key: "TOKEN", Start: -1, End: -1}})
	if err == nil {
		t.Fatal("expected error for negative start")
	}
}

func TestNew_EndBeforeStart_ReturnsError(t *testing.T) {
	_, err := secrettrim.New([]secrettrim.Rule{{Key: "TOKEN", Start: 5, End: 3}})
	if err == nil {
		t.Fatal("expected error for end before start")
	}
}

func TestNew_EmptyKey_ReturnsError(t *testing.T) {
	_, err := secrettrim.New([]secrettrim.Rule{{Key: "", Start: 0, End: -1}})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	tr, _ := secrettrim.New([]secrettrim.Rule{{Key: "TOKEN", Start: 0, End: -1}})
	_, err := tr.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_NoMatchingKey_PassesThrough(t *testing.T) {
	tr, _ := secrettrim.New([]secrettrim.Rule{{Key: "MISSING", Start: 0, End: 5}})
	out, err := tr.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "Bearer abc123xyz" {
		t.Errorf("expected unchanged TOKEN, got %q", out["TOKEN"])
	}
}

func TestApply_TrimFromStart(t *testing.T) {
	tr, _ := secrettrim.New([]secrettrim.Rule{{Key: "TOKEN", Start: 7, End: -1}})
	out, err := tr.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "abc123xyz" {
		t.Errorf("expected %q, got %q", "abc123xyz", out["TOKEN"])
	}
}

func TestApply_TrimBothEnds(t *testing.T) {
	tr, _ := secrettrim.New([]secrettrim.Rule{{Key: "TOKEN", Start: 7, End: 13}})
	out, err := tr.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "abc123" {
		t.Errorf("expected %q, got %q", "abc123", out["TOKEN"])
	}
}

func TestApply_StartBeyondLength_ReturnsEmpty(t *testing.T) {
	tr, _ := secrettrim.New([]secrettrim.Rule{{Key: "SHORT", Start: 100, End: -1}})
	out, err := tr.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SHORT"] != "" {
		t.Errorf("expected empty string, got %q", out["SHORT"])
	}
}

func TestApply_EndBeyondLength_ClampsToEnd(t *testing.T) {
	tr, _ := secrettrim.New([]secrettrim.Rule{{Key: "SHORT", Start: 0, End: 999}})
	out, err := tr.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SHORT"] != "hi" {
		t.Errorf("expected %q, got %q", "hi", out["SHORT"])
	}
}
