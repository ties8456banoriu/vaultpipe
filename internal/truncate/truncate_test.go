package truncate_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/truncate"
)

func TestNew_InvalidMaxLen(t *testing.T) {
	_, err := truncate.New(0)
	if err != truncate.ErrInvalidMaxLen {
		t.Fatalf("expected ErrInvalidMaxLen, got %v", err)
	}
}

func TestNew_ValidMaxLen(t *testing.T) {
	tr, err := truncate.New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	tr, _ := truncate.New(10)
	_, err := tr.Apply(map[string]string{})
	if err != truncate.ErrEmptySecrets {
		t.Fatalf("expected ErrEmptySecrets, got %v", err)
	}
}

func TestApply_ShortValues_Unchanged(t *testing.T) {
	tr, _ := truncate.New(20)
	secrets := map[string]string{"KEY": "short"}
	out, err := tr.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "short" {
		t.Errorf("expected 'short', got %q", out["KEY"])
	}
}

func TestApply_LongValue_Truncated(t *testing.T) {
	tr, _ := truncate.New(5)
	secrets := map[string]string{"TOKEN": "abcdefghij"}
	out, err := tr.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "abcde" {
		t.Errorf("expected 'abcde', got %q", out["TOKEN"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	tr, _ := truncate.New(3)
	secrets := map[string]string{"K": "hello"}
	_, err := tr.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["K"] != "hello" {
		t.Error("original map was mutated")
	}
}

func TestApply_MultipleKeys(t *testing.T) {
	tr, _ := truncate.New(4)
	secrets := map[string]string{
		"A": "123456",
		"B": "xy",
		"C": "abcd",
	}
	out, err := tr.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1234" {
		t.Errorf("A: expected '1234', got %q", out["A"])
	}
	if out["B"] != "xy" {
		t.Errorf("B: expected 'xy', got %q", out["B"])
	}
	if out["C"] != "abcd" {
		t.Errorf("C: expected 'abcd', got %q", out["C"])
	}
}
