package labels_test

import (
	"testing"

	"github.com/fuskovic/vaultpipe/internal/labels"
)

func TestTag_And_Get_RoundTrip(t *testing.T) {
	s := labels.New()
	err := s.Tag("DB_PASSWORD", labels.Set{"env": "prod", "team": "backend"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	set, err := s.Get("DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if set["env"] != "prod" || set["team"] != "backend" {
		t.Errorf("unexpected labels: %v", set)
	}
}

func TestTag_EmptyKey_ReturnsError(t *testing.T) {
	s := labels.New()
	if err := s.Tag("", labels.Set{"env": "prod"}); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestTag_EmptyLabels_ReturnsError(t *testing.T) {
	s := labels.New()
	if err := s.Tag("DB_PASSWORD", labels.Set{}); err == nil {
		t.Fatal("expected error for empty labels")
	}
}

func TestGet_UnknownKey_ReturnsError(t *testing.T) {
	s := labels.New()
	if _, err := s.Get("UNKNOWN"); err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	s := labels.New()
	_ = s.Tag("API_KEY", labels.Set{"env": "dev"})
	set, _ := s.Get("API_KEY")
	set["env"] = "mutated"
	original, _ := s.Get("API_KEY")
	if original["env"] != "dev" {
		t.Error("Get should return a copy, original was mutated")
	}
}

func TestFilter_MatchesSelector(t *testing.T) {
	s := labels.New()
	_ = s.Tag("DB_PASSWORD", labels.Set{"env": "prod", "team": "backend"})
	_ = s.Tag("API_KEY", labels.Set{"env": "dev", "team": "frontend"})
	_ = s.Tag("SECRET_TOKEN", labels.Set{"env": "prod", "team": "frontend"})

	keys := s.Filter(labels.Set{"env": "prod"})
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestFilter_NoMatches_ReturnsEmpty(t *testing.T) {
	s := labels.New()
	_ = s.Tag("DB_PASSWORD", labels.Set{"env": "dev"})
	keys := s.Filter(labels.Set{"env": "prod"})
	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules := []string{"DB_PASSWORD=env=prod,team=backend"}
	out, err := labels.ParseRules(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASSWORD"]["env"] != "prod" {
		t.Errorf("unexpected parse result: %v", out)
	}
}

func TestParseRules_Invalid_ReturnsError(t *testing.T) {
	if _, err := labels.ParseRules([]string{"NOEQUALSSIGN"}); err == nil {
		t.Fatal("expected error for invalid rule")
	}
}
