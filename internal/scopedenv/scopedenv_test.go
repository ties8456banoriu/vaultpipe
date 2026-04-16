package scopedenv_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/scopedenv"
)

func TestNew_EmptyScope_ReturnsError(t *testing.T) {
	_, err := scopedenv.New("")
	if err != scopedenv.ErrNoScope {
		t.Fatalf("expected ErrNoScope, got %v", err)
	}
}

func TestNew_ValidScope_ReturnsScoper(t *testing.T) {
	s, err := scopedenv.New("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Scope() != "dev" {
		t.Fatalf("expected scope dev, got %s", s.Scope())
	}
}

func TestTag_EmptySecrets_ReturnsError(t *testing.T) {
	s, _ := scopedenv.New("dev")
	_, err := s.Tag(map[string]string{})
	if err != scopedenv.ErrEmptySecrets {
		t.Fatalf("expected ErrEmptySecrets, got %v", err)
	}
}

func TestTag_PrefixesKeys(t *testing.T) {
	s, _ := scopedenv.New("staging")
	out, err := s.Tag(map[string]string{"DB_PASS": "secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := out["STAGING_DB_PASS"]; !ok || v != "secret" {
		t.Fatalf("expected STAGING_DB_PASS=secret, got %v", out)
	}
}

func TestFilter_EmptySecrets_ReturnsError(t *testing.T) {
	s, _ := scopedenv.New("dev")
	_, err := s.Filter(map[string]string{})
	if err != scopedenv.ErrEmptySecrets {
		t.Fatalf("expected ErrEmptySecrets, got %v", err)
	}
}

func TestFilter_NoMatchingKeys_ReturnsError(t *testing.T) {
	s, _ := scopedenv.New("prod")
	_, err := s.Filter(map[string]string{"DEV_KEY": "val"})
	if err == nil {
		t.Fatal("expected error for no matching keys")
	}
}

func TestFilter_MatchingKeys_StripsPrefix(t *testing.T) {
	s, _ := scopedenv.New("dev")
	secrets := map[string]string{
		"DEV_API_KEY": "abc",
		"PROD_API_KEY": "xyz",
	}
	out, err := s.Filter(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if out["API_KEY"] != "abc" {
		t.Fatalf("expected API_KEY=abc, got %v", out)
	}
}

func TestNew_TrimsAndLowercasesScope(t *testing.T) {
	s, err := scopedenv.New("  DEV  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Scope() != "dev" {
		t.Fatalf("expected dev, got %s", s.Scope())
	}
}
