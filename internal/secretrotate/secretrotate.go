package secretrotate

import (
	"errors"
	"sync"
	"time"
)

// Policy controls how rotation is triggered.
type Policy string

const (
	PolicyManual    Policy = "manual"
	PolicyScheduled Policy = "scheduled"
)

// RotationRecord holds metadata about a secret rotation event.
type RotationRecord struct {
	EnvKey    string
	VaultPath string
	RotatedAt time.Time
	Policy    Policy
	Version   int
}

// Rotator tracks and enforces secret rotation records.
type Rotator struct {
	mu      sync.RWMutex
	records map[string][]RotationRecord
}

// New returns a new Rotator.
func New() *Rotator {
	return &Rotator{records: make(map[string][]RotationRecord)}
}

// Record adds a rotation event for the given env key.
func (r *Rotator) Record(envKey, vaultPath string, version int, policy Policy) error {
	if envKey == "" {
		return errors.New("secretrotate: env key must not be empty")
	}
	if vaultPath == "" {
		return errors.New("secretrotate: vault path must not be empty")
	}
	if version <= 0 {
		return errors.New("secretrotate: version must be positive")
	}
	rec := RotationRecord{
		EnvKey:    envKey,
		VaultPath: vaultPath,
		RotatedAt: time.Now().UTC(),
		Policy:    policy,
		Version:   version,
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[envKey] = append(r.records[envKey], rec)
	return nil
}

// Latest returns the most recent rotation record for the given key.
func (r *Rotator) Latest(envKey string) (RotationRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list, ok := r.records[envKey]
	if !ok || len(list) == 0 {
		return RotationRecord{}, errors.New("secretrotate: no rotation record found for key")
	}
	return list[len(list)-1], nil
}

// All returns a copy of all rotation records.
func (r *Rotator) All() map[string][]RotationRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string][]RotationRecord, len(r.records))
	for k, v := range r.records {
		cp := make([]RotationRecord, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}

// NeedsRotation returns true if the latest record for the key is older than maxAge.
func (r *Rotator) NeedsRotation(envKey string, maxAge time.Duration) (bool, error) {
	rec, err := r.Latest(envKey)
	if err != nil {
		return false, err
	}
	return time.Since(rec.RotatedAt) > maxAge, nil
}
