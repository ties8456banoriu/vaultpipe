package secretgroup_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretgroup"
)

func TestAdd_And_Get_RoundTrip(t *testing.T) {
	s := secretgroup.New()
	if err := s.Add("db", []string{"DB_HOST", "DB_PASS"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys, err := s.Get("db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 2 || keys[0] != "DB_HOST" || keys[1] != "DB_PASS" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestAdd_EmptyName_ReturnsError(t *testing.T) {
	s := secretgroup.New()
	err := s.Add("", []string{"KEY"})
	if !errors.Is(err, secretgroup.ErrEmptyGroupName) {
		t.Errorf("expected ErrEmptyGroupName, got %v", err)
	}
}

func TestAdd_EmptyKeys_ReturnsError(t *testing.T) {
	s := secretgroup.New()
	err := s.Add("grp", []string{})
	if !errors.Is(err, secretgroup.ErrEmptyKeys) {
		t.Errorf("expected ErrEmptyKeys, got %v", err)
	}
}

func TestGet_NotFound_ReturnsError(t *testing.T) {
	s := secretgroup.New()
	_, err := s.Get("missing")
	if !errors.Is(err, secretgroup.ErrGroupNotFound) {
		t.Errorf("expected ErrGroupNotFound, got %v", err)
	}
}

func TestFilter_ReturnsMatchingSecrets(t *testing.T) {
	s := secretgroup.New()
	_ = s.Add("cache", []string{"REDIS_HOST", "REDIS_PORT"})
	secrets := map[string]string{
		"REDIS_HOST": "localhost",
		"REDIS_PORT": "6379",
		"DB_HOST":    "pghost",
	}
	result, err := s.Filter("cache", secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(result))
	}
	if result["DB_HOST"] != "" {
		t.Errorf("DB_HOST should not be in result")
	}
}

func TestFilter_UnknownGroup_ReturnsError(t *testing.T) {
	s := secretgroup.New()
	_, err := s.Filter("nope", map[string]string{"K": "V"})
	if !errors.Is(err, secretgroup.ErrGroupNotFound) {
		t.Errorf("expected ErrGroupNotFound, got %v", err)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := secretgroup.New()
	_ = s.Add("g1", []string{"A"})
	_ = s.Add("g2", []string{"B", "C"})
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 groups, got %d", len(all))
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	s := secretgroup.New()
	_ = s.Add("g", []string{"X"})
	keys, _ := s.Get("g")
	keys[0] = "MUTATED"
	keys2, _ := s.Get("g")
	if keys2[0] == "MUTATED" {
		t.Error("Get should return a copy")
	}
}
