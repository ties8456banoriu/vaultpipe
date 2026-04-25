package secretpromote_test

import (
	"testing"

	"github.com/elizabethadegbaju/vaultpipe/internal/secretpromote"
)

func TestNew_EmptySource_ReturnsError(t *testing.T) {
	_, err := secretpromote.New("", "prod")
	if err != secretpromote.ErrEmptySource {
		t.Fatalf("expected ErrEmptySource, got %v", err)
	}
}

func TestNew_EmptyTarget_ReturnsError(t *testing.T) {
	_, err := secretpromote.New("staging", "")
	if err != secretpromote.ErrEmptyTarget {
		t.Fatalf("expected ErrEmptyTarget, got %v", err)
	}
}

func TestNew_SameNamespace_ReturnsError(t *testing.T) {
	_, err := secretpromote.New("staging", "staging")
	if err != secretpromote.ErrSameNamespace {
		t.Fatalf("expected ErrSameNamespace, got %v", err)
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	p, err := secretpromote.New("staging", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = p.Apply(map[string]string{})
	if err != secretpromote.ErrEmptySecrets {
		t.Fatalf("expected ErrEmptySecrets, got %v", err)
	}
}

func TestApply_PromotesMatchingKeys(t *testing.T) {
	p, _ := secretpromote.New("staging", "prod")
	secrets := map[string]string{
		"staging_DB_URL": "postgres://staging",
		"staging_API_KEY": "key-abc",
		"other_VAR": "should-be-ignored",
	}
	result, err := p.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Promoted["prod_DB_URL"] != "postgres://staging" {
		t.Errorf("expected prod_DB_URL to be promoted")
	}
	if result.Promoted["prod_API_KEY"] != "key-abc" {
		t.Errorf("expected prod_API_KEY to be promoted")
	}
	if _, ok := result.Promoted["other_VAR"]; ok {
		t.Errorf("expected other_VAR to be ignored")
	}
}

func TestApply_SkipsExistingTargetWithoutOverwrite(t *testing.T) {
	p, _ := secretpromote.New("staging", "prod")
	secrets := map[string]string{
		"staging_DB_URL": "postgres://staging",
		"prod_DB_URL":    "postgres://prod",
	}
	result, err := p.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) == 0 {
		t.Errorf("expected prod_DB_URL to be skipped")
	}
	if _, ok := result.Promoted["prod_DB_URL"]; ok {
		t.Errorf("expected prod_DB_URL not to be overwritten")
	}
}

func TestApply_OverwritesExistingTargetWhenEnabled(t *testing.T) {
	p, _ := secretpromote.New("staging", "prod", secretpromote.WithOverwrite())
	secrets := map[string]string{
		"staging_DB_URL": "postgres://new",
		"prod_DB_URL":    "postgres://old",
	}
	result, err := p.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Promoted["prod_DB_URL"] != "postgres://new" {
		t.Errorf("expected prod_DB_URL to be overwritten, got %q", result.Promoted["prod_DB_URL"])
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected no skipped keys, got %v", result.Skipped)
	}
}

func TestApply_NoMatchingKeys_ReturnsEmptyPromoted(t *testing.T) {
	p, _ := secretpromote.New("staging", "prod")
	secrets := map[string]string{
		"dev_KEY": "value",
	}
	result, err := p.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Promoted) != 0 {
		t.Errorf("expected no promoted keys, got %v", result.Promoted)
	}
}
