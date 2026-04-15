package redact_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/redact"
)

func TestMask_EmptyValue(t *testing.T) {
	r := redact.NewRedactor()
	if got := r.Mask(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestMask_ShortValue_FullyMasked(t *testing.T) {
	r := redact.NewRedactor()
	got := r.Mask("abc")
	if got != "***" {
		t.Errorf("expected '***', got %q", got)
	}
}

func TestMask_LongValue_TrailingVisible(t *testing.T) {
	r := redact.NewRedactor()
	value := "supersecretvalue"
	got := r.Mask(value)

	if !strings.HasSuffix(got, "alue") {
		t.Errorf("expected suffix 'alue', got %q", got)
	}
	if len(got) != len(value) {
		t.Errorf("expected length %d, got %d", len(value), len(got))
	}
	maskedPart := got[:len(got)-redact.DefaultVisibleChars]
	for _, ch := range maskedPart {
		if string(ch) != redact.DefaultMaskChar {
			t.Errorf("expected mask char in prefix, got %q", string(ch))
		}
	}
}

func TestMaskAll_ReplacesAll(t *testing.T) {
	r := redact.NewRedactor()
	got := r.MaskAll("topsecret")
	if got != "*********" {
		t.Errorf("expected all stars, got %q", got)
	}
}

func TestMaskAll_Empty(t *testing.T) {
	r := redact.NewRedactor()
	if got := r.MaskAll(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestMaskMap_MasksValues_KeepsKeys(t *testing.T) {
	r := redact.NewRedactor()
	secrets := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abcdefghij",
	}
	masked := r.MaskMap(secrets)

	if _, ok := masked["DB_PASSWORD"]; !ok {
		t.Error("expected key DB_PASSWORD to be present")
	}
	if masked["DB_PASSWORD"] == "hunter2" {
		t.Error("expected DB_PASSWORD value to be masked")
	}
	if masked["API_KEY"] == "abcdefghij" {
		t.Error("expected API_KEY value to be masked")
	}
	if !strings.HasSuffix(masked["API_KEY"], "ghij") {
		t.Errorf("expected API_KEY masked value to end with 'ghij', got %q", masked["API_KEY"])
	}
}

func TestMaskMap_DoesNotMutateOriginal(t *testing.T) {
	r := redact.NewRedactor()
	secrets := map[string]string{"TOKEN": "plaintext"}
	_ = r.MaskMap(secrets)
	if secrets["TOKEN"] != "plaintext" {
		t.Error("original map should not be mutated")
	}
}
