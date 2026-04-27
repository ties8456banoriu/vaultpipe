package secretsearch_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretsearch"
)

var base = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "s3cr3t",
	"API_KEY":     "abc123",
	"API_SECRET":  "xyz789",
}

func TestNew_UnsupportedMode_ReturnsError(t *testing.T) {
	_, err := secretsearch.New("fuzzy", true, false)
	if err == nil {
		t.Fatal("expected error for unsupported mode")
	}
}

func TestNew_NeitherFieldEnabled_ReturnsError(t *testing.T) {
	_, err := secretsearch.New(secretsearch.ModeExact, false, false)
	if err == nil {
		t.Fatal("expected error when both searchKey and searchVal are false")
	}
}

func TestSearch_EmptySecrets_ReturnsError(t *testing.T) {
	s, _ := secretsearch.New(secretsearch.ModeExact, true, false)
	_, err := s.Search(map[string]string{}, "DB_HOST")
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestSearch_EmptyQuery_ReturnsError(t *testing.T) {
	s, _ := secretsearch.New(secretsearch.ModeExact, true, false)
	_, err := s.Search(base, "")
	if err == nil {
		t.Fatal("expected error for empty query")
	}
}

func TestSearch_ExactKey_ReturnsMatch(t *testing.T) {
	s, _ := secretsearch.New(secretsearch.ModeExact, true, false)
	results, err := s.Search(base, "DB_HOST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "DB_HOST" {
		t.Fatalf("expected 1 result for DB_HOST, got %v", results)
	}
}

func TestSearch_PrefixKey_ReturnsMultiple(t *testing.T) {
	s, _ := secretsearch.New(secretsearch.ModePrefix, true, false)
	results, err := s.Search(base, "API_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results for prefix API_, got %d", len(results))
	}
}

func TestSearch_RegexValue_ReturnsMatch(t *testing.T) {
	s, _ := secretsearch.New(secretsearch.ModeRegex, false, true)
	results, err := s.Search(base, `^[a-z0-9]+$`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one match for lowercase regex")
	}
}

func TestSearch_InvalidRegex_ReturnsError(t *testing.T) {
	s, _ := secretsearch.New(secretsearch.ModeRegex, true, false)
	_, err := s.Search(base, "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestSearch_NoMatch_ReturnsEmptySlice(t *testing.T) {
	s, _ := secretsearch.New(secretsearch.ModeExact, true, true)
	results, err := s.Search(base, "NONEXISTENT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no results, got %v", results)
	}
}
