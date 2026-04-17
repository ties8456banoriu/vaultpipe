package secretmeta_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/secretmeta"
)

func TestRecord_And_Get_RoundTrip(t *testing.T) {
	s := secretmeta.New()
	m := secretmeta.Meta{
		EnvKey:    "DB_PASSWORD",
		VaultPath: "secret/data/db",
		Mount:     "secret",
		Version:   3,
	}
	if err := s.Record(m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.Get("DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.VaultPath != m.VaultPath {
		t.Errorf("expected %q got %q", m.VaultPath, got.VaultPath)
	}
	if got.Version != 3 {
		t.Errorf("expected version 3 got %d", got.Version)
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	s := secretmeta.New()
	before := time.Now().UTC()
	_ = s.Record(secretmeta.Meta{EnvKey: "X", VaultPath: "secret/data/x"})
	got, _ := s.Get("X")
	if got.FetchedAt.Before(before) {
		t.Error("expected FetchedAt to be set to now")
	}
}

func TestRecord_EmptyEnvKey_ReturnsError(t *testing.T) {
	s := secretmeta.New()
	err := s.Record(secretmeta.Meta{VaultPath: "secret/data/x"})
	if err == nil {
		t.Fatal("expected error for empty env key")
	}
}

func TestRecord_EmptyVaultPath_ReturnsError(t *testing.T) {
	s := secretmeta.New()
	err := s.Record(secretmeta.Meta{EnvKey: "X"})
	if err == nil {
		t.Fatal("expected error for empty vault path")
	}
}

func TestGet_UnknownKey_ReturnsError(t *testing.T) {
	s := secretmeta.New()
	_, err := s.Get("UNKNOWN")
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := secretmeta.New()
	_ = s.Record(secretmeta.Meta{EnvKey: "A", VaultPath: "secret/data/a"})
	_ = s.Record(secretmeta.Meta{EnvKey: "B", VaultPath: "secret/data/b"})
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 records got %d", len(all))
	}
}

func TestDelete_RemovesRecord(t *testing.T) {
	s := secretmeta.New()
	_ = s.Record(secretmeta.Meta{EnvKey: "A", VaultPath: "secret/data/a"})
	s.Delete("A")
	_, err := s.Get("A")
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}
