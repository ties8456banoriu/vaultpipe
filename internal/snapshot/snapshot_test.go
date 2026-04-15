package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/snapshot"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshot.json")
}

func TestStore_And_Load_RoundTrip(t *testing.T) {
	path := tempPath(t)
	secrets := map[string]string{"DB_PASS": "s3cr3t", "API_KEY": "abc123"}

	if err := snapshot.Store(path, secrets); err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if snap.Secrets["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected DB_PASS=s3cr3t, got %s", snap.Secrets["DB_PASS"])
	}
	if snap.Secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %s", snap.Secrets["API_KEY"])
	}
}

func TestStore_SetsTimestamp(t *testing.T) {
	path := tempPath(t)
	before := time.Now().UTC().Add(-time.Second)

	if err := snapshot.Store(path, map[string]string{"KEY": "val"}); err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	snap, _ := snapshot.Load(path)
	if snap.CapturedAt.Before(before) {
		t.Error("expected CapturedAt to be recent")
	}
}

func TestStore_EmptySecrets_ReturnsError(t *testing.T) {
	path := tempPath(t)
	err := snapshot.Store(path, map[string]string{})
	if err == nil {
		t.Error("expected error for empty secrets, got nil")
	}
}

func TestLoad_NoFile_ReturnsErrNoSnapshot(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err != snapshot.ErrNoSnapshot {
		t.Errorf("expected ErrNoSnapshot, got %v", err)
	}
}

func TestStore_FilePermissions(t *testing.T) {
	path := tempPath(t)
	if err := snapshot.Store(path, map[string]string{"X": "y"}); err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file mode 0600, got %v", info.Mode().Perm())
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not json{"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
