// Package rotate provides utilities for detecting and handling secret
// rotation events by comparing current secrets against a stored baseline.
package rotate

import (
	"errors"
	"fmt"
	"time"

	"github.com/yourusername/vaultpipe/internal/diff"
)

// ErrNoBaseline is returned when no baseline snapshot exists to compare against.
var ErrNoBaseline = errors.New("rotate: no baseline set")

// RotationEvent describes a detected secret rotation.
type RotationEvent struct {
	DetectedAt time.Time
	Changes    []diff.Change
}

// Detector tracks a baseline set of secrets and detects when they rotate.
type Detector struct {
	baseline map[string]string
}

// NewDetector creates a new Detector with no baseline set.
func NewDetector() *Detector {
	return &Detector{}
}

// SetBaseline records the current secrets as the baseline for future comparisons.
func (d *Detector) SetBaseline(secrets map[string]string) error {
	if len(secrets) == 0 {
		return errors.New("rotate: cannot set empty baseline")
	}
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	d.baseline = copy
	return nil
}

// Detect compares the provided secrets against the baseline.
// Returns a RotationEvent if any changes are found, or ErrNoBaseline if
// SetBaseline has not been called yet.
func (d *Detector) Detect(current map[string]string) (*RotationEvent, error) {
	if d.baseline == nil {
		return nil, ErrNoBaseline
	}
	changes := diff.Compare(d.baseline, current)
	if len(changes) == 0 {
		return nil, nil
	}
	return &RotationEvent{
		DetectedAt: time.Now().UTC(),
		Changes:    changes,
	}, nil
}

// Summary returns a human-readable summary of the rotation event.
func (e *RotationEvent) Summary() string {
	return fmt.Sprintf("rotation detected at %s: %d key(s) changed",
		e.DetectedAt.Format(time.RFC3339), len(e.Changes))
}
