package env_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/env"
)

func TestResolve_NoMatchingVars_ReturnsErrNoOverrides(t *testing.T) {
	r := env.NewResolver("VAULTPIPE_OVERRIDE_ZZZNOMATCH_")
	_, err := r.Resolve()
	if err != env.ErrNoOverrides {
		t.Fatalf("expected ErrNoOverrides, got %v", err)
	}
}

func TestResolve_MatchingVars_ReturnsMap(t *testing.T) {
	t.Setenv("VAULTPIPE_OVERRIDE_DB_HOST", "localhost")
	t.Setenv("VAULTPIPE_OVERRIDE_DB_PORT", "5432")

	r := env.NewResolver("VAULTPIPE_OVERRIDE_")
	got, err := r.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: want %q, got %q", "localhost", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT: want %q, got %q", "5432", got["DB_PORT"])
	}
}

func TestResolve_PrefixStrippedAndUppercased(t *testing.T) {
	t.Setenv("VAULTPIPE_OVERRIDE_my_key", "value1")

	r := env.NewResolver("vaultpipe_override_") // lowercase prefix
	got, err := r.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := got["MY_KEY"]; !ok {
		t.Errorf("expected MY_KEY in result, got keys: %v", got)
	}
}

func TestApply_OverridesReplaceBase(t *testing.T) {
	base := map[string]string{
		"DB_HOST": "vault-host",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}
	overrides := map[string]string{
		"DB_HOST": "localhost",
	}

	got := env.Apply(base, overrides)

	if got["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: want %q, got %q", "localhost", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT should be preserved, got %q", got["DB_PORT"])
	}
	if got["API_KEY"] != "secret" {
		t.Errorf("API_KEY should be preserved, got %q", got["API_KEY"])
	}
}

func TestApply_EmptyOverrides_ReturnsBaseUnchanged(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	got := env.Apply(base, map[string]string{})
	if got["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", got["FOO"])
	}
}

func TestApply_DoesNotMutateBase(t *testing.T) {
	base := map[string]string{"FOO": "original"}
	overrides := map[string]string{"FOO": "changed"}
	env.Apply(base, overrides)
	if base["FOO"] != "original" {
		t.Error("Apply must not mutate the base map")
	}
}
