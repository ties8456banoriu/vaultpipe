package secretexpand_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretexpand"
)

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	e := secretexpand.New()
	_, err := e.Apply(map[string]string{})
	if !errors.Is(err, secretexpand.ErrEmptySecrets) {
		t.Fatalf("expected ErrEmptySecrets, got %v", err)
	}
}

func TestApply_NoRefs_ReturnsUnchanged(t *testing.T) {
	e := secretexpand.New()
	input := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out, err := e.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" || out["DB_PORT"] != "5432" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestApply_SimpleRef_Resolved(t *testing.T) {
	e := secretexpand.New()
	input := map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://${HOST}/mydb",
	}
	out, err := e.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost/mydb" {
		t.Fatalf("expected resolved DSN, got %q", out["DSN"])
	}
}

func TestApply_ChainedRefs_Resolved(t *testing.T) {
	e := secretexpand.New()
	input := map[string]string{
		"A": "hello",
		"B": "${A}_world",
		"C": "${B}!",
	}
	out, err := e.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["C"] != "hello_world!" {
		t.Fatalf("expected 'hello_world!', got %q", out["C"])
	}
}

func TestApply_UnknownRef_ReturnsError(t *testing.T) {
	e := secretexpand.New()
	input := map[string]string{
		"DSN": "postgres://${MISSING_HOST}/db",
	}
	_, err := e.Apply(input)
	if err == nil {
		t.Fatal("expected error for unknown reference")
	}
}

func TestApply_CircularRef_ReturnsError(t *testing.T) {
	e := secretexpand.New()
	input := map[string]string{
		"A": "${B}",
		"B": "${A}",
	}
	_, err := e.Apply(input)
	if !errors.Is(err, secretexpand.ErrCircularReference) {
		t.Fatalf("expected ErrCircularReference, got %v", err)
	}
}

func TestApply_MultipleRefsInValue(t *testing.T) {
	e := secretexpand.New()
	input := map[string]string{
		"USER": "admin",
		"PASS": "secret",
		"URL":  "${USER}:${PASS}@host",
	}
	out, err := e.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "admin:secret@host" {
		t.Fatalf("expected resolved URL, got %q", out["URL"])
	}
}
