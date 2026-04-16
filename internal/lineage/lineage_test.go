package lineage_test

import (
	"testing"

	"github.com/foulkelore/vaultpipe/internal/lineage"
)

func TestTrack_And_Get_RoundTrip(t *testing.T) {
	tr := lineage.New()
	if err := tr.Track("DB_HOST", "secret/data/app", "db_host"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	o, err := tr.Get("DB_HOST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if o.VaultPath != "secret/data/app" || o.VaultKey != "db_host" || o.EnvKey != "DB_HOST" {
		t.Errorf("unexpected origin: %+v", o)
	}
}

func TestTrack_EmptyEnvKey_ReturnsError(t *testing.T) {
	tr := lineage.New()
	if err := tr.Track("", "secret/data/app", "key"); err == nil {
		t.Fatal("expected error for empty envKey")
	}
}

func TestTrack_EmptyVaultPath_ReturnsError(t *testing.T) {
	tr := lineage.New()
	if err := tr.Track("KEY", "", "key"); err == nil {
		t.Fatal("expected error for empty vaultPath")
	}
}

func TestGet_UnknownKey_ReturnsError(t *testing.T) {
	tr := lineage.New()
	if _, err := tr.Get("MISSING"); err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := lineage.New()
	_ = tr.Track("A", "secret/data/x", "a")
	_ = tr.Track("B", "secret/data/x", "b")
	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 records, got %d", len(all))
	}
	// Mutating the copy must not affect the tracker.
	delete(all, "A")
	if _, err := tr.Get("A"); err != nil {
		t.Error("tracker record was mutated via All() copy")
	}
}

func TestTrackAll_PopulatesRecords(t *testing.T) {
	tr := lineage.New()
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := tr.TrackAll(secrets, "secret/data/svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k := range secrets {
		o, err := tr.Get(k)
		if err != nil {
			t.Errorf("missing record for %q: %v", k, err)
		}
		if o.VaultPath != "secret/data/svc" {
			t.Errorf("wrong vault path for %q: %s", k, o.VaultPath)
		}
	}
}

func TestTrackAll_EmptySecrets_ReturnsError(t *testing.T) {
	tr := lineage.New()
	if err := tr.TrackAll(map[string]string{}, "secret/data/svc"); err == nil {
		t.Fatal("expected error for empty secrets")
	}
}
