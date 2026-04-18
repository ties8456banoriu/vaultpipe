package secretalias_test

import (
	"testing"

	"github.com/bjarnemagnussen/vaultpipe/internal/secretalias"
)

func TestAdd_And_Aliases_RoundTrip(t *testing.T) {
	a := secretalias.New()
	if err := a.Add("DB_PASSWORD", "DATABASE_PASSWORD", "DB_PASS"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	aliases, err := a.Aliases("DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(aliases) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(aliases))
	}
}

func TestAdd_EmptyOriginal_ReturnsError(t *testing.T) {
	a := secretalias.New()
	if err := a.Add("", "ALIAS"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdd_EmptyAlias_ReturnsError(t *testing.T) {
	a := secretalias.New()
	if err := a.Add("KEY", ""); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdd_DuplicateAliasConflict_ReturnsError(t *testing.T) {
	a := secretalias.New()
	_ = a.Add("KEY_A", "SHARED")
	if err := a.Add("KEY_B", "SHARED"); err == nil {
		t.Fatal("expected conflict error, got nil")
	}
}

func TestResolve_ReturnsOriginal(t *testing.T) {
	a := secretalias.New()
	_ = a.Add("API_KEY", "SERVICE_API_KEY")
	orig, err := a.Resolve("SERVICE_API_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orig != "API_KEY" {
		t.Fatalf("expected API_KEY, got %s", orig)
	}
}

func TestResolve_UnknownAlias_ReturnsErrNoAlias(t *testing.T) {
	a := secretalias.New()
	_, err := a.Resolve("UNKNOWN")
	if err != secretalias.ErrNoAlias {
		t.Fatalf("expected ErrNoAlias, got %v", err)
	}
}

func TestAliases_UnknownKey_ReturnsErrNoAlias(t *testing.T) {
	a := secretalias.New()
	_, err := a.Aliases("MISSING")
	if err != secretalias.ErrNoAlias {
		t.Fatalf("expected ErrNoAlias, got %v", err)
	}
}

func TestApply_ExpandsAliases(t *testing.T) {
	a := secretalias.New()
	_ = a.Add("DB_PASSWORD", "DATABASE_PASSWORD")
	secrets := map[string]string{"DB_PASSWORD": "s3cr3t"}
	out, err := a.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATABASE_PASSWORD"] != "s3cr3t" {
		t.Fatalf("expected alias to be expanded, got %q", out["DATABASE_PASSWORD"])
	}
	if out["DB_PASSWORD"] != "s3cr3t" {
		t.Fatalf("original key should be preserved")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	a := secretalias.New()
	_, err := a.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
