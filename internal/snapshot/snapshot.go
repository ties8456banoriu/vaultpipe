// Package snapshot provides functionality for persisting and loading
// secret snapshots to disk, enabling offline fallback and change detection.
package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// ErrNoSnapshot is returned when no snapshot file exists at the given path.
var ErrNoSnapshot = errors.New("no snapshot found")

// Snapshot holds a point-in-time capture of fetched secrets.
type Snapshot struct {
	CapturedAt time.Time         `json:"captured_at"`
	Secrets    map[string]string `json:"secrets"`
}

// Store persists a secrets map to a JSON snapshot file at path.
func Store(path string, secrets map[string]string) error {
	if len(secrets) == 0 {
		return errors.New("cannot store empty secrets snapshot")
	}

	snap := Snapshot{
		CapturedAt: time.Now().UTC(),
		Secrets:    secrets,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// Load reads a snapshot from path and returns it.
// Returns ErrNoSnapshot if the file does not exist.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoSnapshot
		}
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}

	return &snap, nil
}
