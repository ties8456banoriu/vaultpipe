package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpipe/internal/profile"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "profiles.json")
}

func TestSet_And_Get_RoundTrip(t *testing.T) {
	s := profile.NewStore(tempStorePath(t))
	p := profile.Profile{Name: "dev", SecretPath: "secret/dev", EnvFile: ".env"}
	if err := s.Set(p); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := s.Get("dev")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.SecretPath != p.SecretPath {
		t.Errorf("SecretPath: got %q want %q", got.SecretPath, p.SecretPath)
	}
}

func TestSet_EmptyName_ReturnsError(t *testing.T) {
	s := profile.NewStore(tempStorePath(t))
	err := s.Set(profile.Profile{SecretPath: "secret/dev"})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestSet_EmptySecretPath_ReturnsError(t *testing.T) {
	s := profile.NewStore(tempStorePath(t))
	err := s.Set(profile.Profile{Name: "dev"})
	if err == nil {
		t.Fatal("expected error for empty secret_path")
	}
}

func TestGet_NotFound_ReturnsError(t *testing.T) {
	s := profile.NewStore(tempStorePath(t))
	_, err := s.Get("missing")
	if !errors.Is(err, profile.ErrProfileNotFound) {
		t.Errorf("expected ErrProfileNotFound, got %v", err)
	}
}

func TestDelete_RemovesProfile(t *testing.T) {
	s := profile.NewStore(tempStorePath(t))
	_ = s.Set(profile.Profile{Name: "dev", SecretPath: "secret/dev"})
	if err := s.Delete("dev"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Get("dev"); err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestDelete_NotFound_ReturnsError(t *testing.T) {
	s := profile.NewStore(tempStorePath(t))
	err := s.Delete("ghost")
	if !errors.Is(err, profile.ErrProfileNotFound) {
		t.Errorf("expected ErrProfileNotFound, got %v", err)
	}
}

func TestSave_And_Load_RoundTrip(t *testing.T) {
	path := tempStorePath(t)
	s1 := profile.NewStore(path)
	_ = s1.Set(profile.Profile{Name: "staging", SecretPath: "secret/staging", EnvFile: ".env.staging"})
	if err := s1.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s2 := profile.NewStore(path)
	if err := s2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	p, err := s2.Get("staging")
	if err != nil {
		t.Fatalf("Get after load: %v", err)
	}
	if p.EnvFile != ".env.staging" {
		t.Errorf("EnvFile: got %q want .env.staging", p.EnvFile)
	}
}

func TestLoad_NoFile_ReturnsErrNoProfiles(t *testing.T) {
	s := profile.NewStore(filepath.Join(t.TempDir(), "nonexistent.json"))
	err := s.Load()
	if !errors.Is(err, profile.ErrNoProfiles) {
		t.Errorf("expected ErrNoProfiles, got %v", err)
	}
}

func TestList_ReturnsAllProfiles(t *testing.T) {
	s := profile.NewStore(tempStorePath(t))
	_ = s.Set(profile.Profile{Name: "dev", SecretPath: "secret/dev"})
	_ = s.Set(profile.Profile{Name: "prod", SecretPath: "secret/prod"})
	list := s.List()
	if len(list) != 2 {
		t.Errorf("List: got %d profiles, want 2", len(list))
	}
}

func init() {
	// ensure errors package is available
	_ = os.Stderr
}
