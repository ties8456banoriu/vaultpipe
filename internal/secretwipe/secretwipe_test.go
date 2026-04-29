package secretwipe_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretwipe"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"APP_NAME":    "myapp",
	}
}

func TestWipe_EmptySecrets_ReturnsError(t *testing.T) {
	w := secretwipe.New()
	err := w.Wipe(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets, got nil")
	}
}

func TestWipe_NoPatterns_WipesAll(t *testing.T) {
	w := secretwipe.New()
	secrets := baseSecrets()
	if err := w.Wipe(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range secrets {
		if v != "" {
			t.Errorf("expected key %q to be wiped, got %q", k, v)
		}
	}
}

func TestWipe_WithPattern_WipesMatchingOnly(t *testing.T) {
	w := secretwipe.New(secretwipe.WithPatterns([]string{"DB_*"}))
	secrets := baseSecrets()
	if err := w.Wipe(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["DB_PASSWORD"] != "" {
		t.Errorf("expected DB_PASSWORD to be wiped")
	}
	if secrets["API_KEY"] == "" {
		t.Errorf("expected API_KEY to be retained")
	}
	if secrets["APP_NAME"] == "" {
		t.Errorf("expected APP_NAME to be retained")
	}
}

func TestWipeCopy_EmptySecrets_ReturnsError(t *testing.T) {
	w := secretwipe.New()
	_, err := w.WipeCopy(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets, got nil")
	}
}

func TestWipeCopy_DoesNotMutateOriginal(t *testing.T) {
	w := secretwipe.New()
	secrets := baseSecrets()
	out, err := w.WipeCopy(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range out {
		if v != "" {
			t.Errorf("expected copy key %q to be wiped", k)
		}
	}
	if secrets["DB_PASSWORD"] == "" {
		t.Error("original map should not be mutated")
	}
}

func TestWipe_ExactPattern_MatchesExact(t *testing.T) {
	w := secretwipe.New(secretwipe.WithPatterns([]string{"API_KEY"}))
	secrets := baseSecrets()
	if err := w.Wipe(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["API_KEY"] != "" {
		t.Errorf("expected API_KEY to be wiped")
	}
	if secrets["DB_PASSWORD"] == "" {
		t.Errorf("expected DB_PASSWORD to be retained")
	}
}
