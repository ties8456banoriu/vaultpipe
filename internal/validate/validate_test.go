package validate_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/validate"
)

func TestValidate_EmptySecrets_ReturnsError(t *testing.T) {
	v := validate.NewValidator(false)
	_, err := v.Validate(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets, got nil")
	}
}

func TestValidate_ValidKeys_NoError(t *testing.T) {
	v := validate.NewValidator(false)
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"_PRIVATE": "value",
		"mixedCase123": "ok",
	}
	warnings, err := v.Validate(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warnings))
	}
}

func TestValidate_InvalidKey_ReturnsError(t *testing.T) {
	v := validate.NewValidator(false)
	secrets := map[string]string{
		"1INVALID": "value",
	}
	_, err := v.Validate(secrets)
	if err == nil {
		t.Fatal("expected error for invalid key name, got nil")
	}
}

func TestValidate_EmptyValue_WarnOnEmpty(t *testing.T) {
	v := validate.NewValidator(true)
	secrets := map[string]string{
		"MY_SECRET": "",
	}
	warnings, err := v.Validate(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Key != "MY_SECRET" {
		t.Errorf("expected warning for MY_SECRET, got %q", warnings[0].Key)
	}
}

func TestValidate_EmptyValue_NoWarnWhenDisabled(t *testing.T) {
	v := validate.NewValidator(false)
	secrets := map[string]string{
		"MY_SECRET": "",
	}
	warnings, err := v.Validate(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warnings))
	}
}

func TestValidate_MixedKeys_ErrorAndWarning(t *testing.T) {
	v := validate.NewValidator(true)
	secrets := map[string]string{
		"GOOD_KEY": "",
		"bad-key":  "value",
	}
	warnings, err := v.Validate(secrets)
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
	// GOOD_KEY should still produce a warning despite the error
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
}
