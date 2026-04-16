// Package lineage tracks the origin of each secret value — which Vault path
// and key it was sourced from — so downstream components can provide
// provenance information in audit logs and diffs.
package lineage

import (
	"errors"
	"fmt"
)

// Origin describes where a single secret value came from.
type Origin struct {
	VaultPath string // e.g. "secret/data/myapp"
	VaultKey  string // e.g. "DB_PASSWORD"
	EnvKey    string // mapped env key, e.g. "DATABASE_PASSWORD"
}

// Record maps each env key to its Origin.
type Record map[string]Origin

// Tracker builds and queries lineage records.
type Tracker struct {
	records Record
}

// New returns an empty Tracker.
func New() *Tracker {
	return &Tracker{records: make(Record)}
}

// Track registers the origin of an env key.
func (t *Tracker) Track(envKey, vaultPath, vaultKey string) error {
	if envKey == "" {
		return errors.New("lineage: envKey must not be empty")
	}
	if vaultPath == "" {
		return errors.New("lineage: vaultPath must not be empty")
	}
	if vaultKey == "" {
		return errors.New("lineage: vaultKey must not be empty")
	}
	t.records[envKey] = Origin{
		VaultPath: vaultPath,
		VaultKey:  vaultKey,
		EnvKey:    envKey,
	}
	return nil
}

// Get returns the Origin for the given env key.
func (t *Tracker) Get(envKey string) (Origin, error) {
	o, ok := t.records[envKey]
	if !ok {
		return Origin{}, fmt.Errorf("lineage: no record for key %q", envKey)
	}
	return o, nil
}

// All returns a copy of all tracked records.
func (t *Tracker) All() Record {
	out := make(Record, len(t.records))
	for k, v := range t.records {
		out[k] = v
	}
	return out
}

// TrackAll registers origins for every key in secrets using a shared vaultPath.
// The vaultKey is assumed to equal the env key before any mapping.
func (t *Tracker) TrackAll(secrets map[string]string, vaultPath string) error {
	if len(secrets) == 0 {
		return errors.New("lineage: secrets must not be empty")
	}
	for k := range secrets {
		if err := t.Track(k, vaultPath, k); err != nil {
			return err
		}
	}
	return nil
}
