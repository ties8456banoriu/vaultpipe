package filter

import (
	"testing"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"API_KEY":     "abc123",
		"APP_DEBUG":   "true",
	}
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	f := NewFilter(nil)
	result, err := f.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 4 {
		t.Errorf("expected 4 keys, got %d", len(result))
	}
}

func TestApply_IncludeWildcard(t *testing.T) {
	rules, _ := ParseRules([]string{"DB_*"})
	f := NewFilter(rules)
	result, err := f.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be included")
	}
}

func TestApply_ExcludeRule(t *testing.T) {
	rules, _ := ParseRules([]string{"DB_*", "!DB_PASSWORD"})
	f := NewFilter(rules)
	result, err := f.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["DB_PASSWORD"]; ok {
		t.Error("expected DB_PASSWORD to be excluded")
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be included")
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	f := NewFilter(nil)
	_, err := f.Apply(map[string]string{})
	if err == nil {
		t.Error("expected error for empty secrets")
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"API_*", "!API_SECRET"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Exclude {
		t.Error("first rule should not be exclude")
	}
	if !rules[1].Exclude {
		t.Error("second rule should be exclude")
	}
}

func TestParseRules_EmptyString_ReturnsError(t *testing.T) {
	_, err := ParseRules([]string{"API_*", ""})
	if err == nil {
		t.Error("expected error for empty rule string")
	}
}

func TestMatchPattern_ExactMatch(t *testing.T) {
	if !matchPattern("DB_HOST", "DB_HOST") {
		t.Error("expected exact match")
	}
	if matchPattern("DB_HOST", "DB_PORT") {
		t.Error("expected no match")
	}
}

func TestMatchPattern_Wildcard(t *testing.T) {
	if !matchPattern("*", "ANYTHING") {
		t.Error("* should match anything")
	}
	if !matchPattern("DB_*", "DB_HOST") {
		t.Error("DB_* should match DB_HOST")
	}
	if matchPattern("DB_*", "API_KEY") {
		t.Error("DB_* should not match API_KEY")
	}
}
