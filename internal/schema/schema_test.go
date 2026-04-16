package schema_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/schema"
)

func TestParseRules_Valid(t *testing.T) {
	rules, err := schema.ParseRules([]string{"DB_HOST:required:string", "PORT:required:int", "DEBUG:bool"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
}

func TestParseRules_UnknownModifier(t *testing.T) {
	_, err := schema.ParseRules([]string{"KEY:unknown"})
	if err == nil {
		t.Fatal("expected error for unknown modifier")
	}
}

func TestParseRules_EmptyKey(t *testing.T) {
	_, err := schema.ParseRules([]string{":required"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	v := schema.NewValidator([]schema.FieldRule{
		{Key: "DB_HOST", Required: true},
	})
	err := v.Validate(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
}

func TestValidate_RequiredPresent(t *testing.T) {
	v := schema.NewValidator([]schema.FieldRule{
		{Key: "DB_HOST", Required: true},
	})
	err := v.Validate(map[string]string{"DB_HOST": "localhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_IntPatternValid(t *testing.T) {
	v := schema.NewValidator([]schema.FieldRule{
		{Key: "PORT", Pattern: "int"},
	})
	if err := v.Validate(map[string]string{"PORT": "5432"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_IntPatternInvalid(t *testing.T) {
	v := schema.NewValidator([]schema.FieldRule{
		{Key: "PORT", Pattern: "int"},
	})
	if err := v.Validate(map[string]string{"PORT": "abc"}); err == nil {
		t.Fatal("expected error for non-int value")
	}
}

func TestValidate_BoolPatternValid(t *testing.T) {
	v := schema.NewValidator([]schema.FieldRule{
		{Key: "DEBUG", Pattern: "bool"},
	})
	for _, val := range []string{"true", "false", "True", "FALSE"} {
		if err := v.Validate(map[string]string{"DEBUG": val}); err != nil {
			t.Fatalf("unexpected error for %q: %v", val, err)
		}
	}
}

func TestValidate_BoolPatternInvalid(t *testing.T) {
	v := schema.NewValidator([]schema.FieldRule{
		{Key: "DEBUG", Pattern: "bool"},
	})
	if err := v.Validate(map[string]string{"DEBUG": "yes"}); err == nil {
		t.Fatal("expected error for invalid bool")
	}
}

func TestValidate_EmptySecrets_NoRules_OK(t *testing.T) {
	v := schema.NewValidator([]schema.FieldRule{})
	if err := v.Validate(map[string]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
