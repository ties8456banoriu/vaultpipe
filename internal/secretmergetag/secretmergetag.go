// Package secretmergetag merges tag maps from multiple secret sets.
package secretmergetag

import (
	"errors"
	"fmt"
)

// ConflictPolicy controls how tag key conflicts are resolved.
type ConflictPolicy string

const (
	PolicySkip      ConflictPolicy = "skip"
	PolicyOverwrite ConflictPolicy = "overwrite"
	PolicyError     ConflictPolicy = "error"
)

// Merger merges tag maps across secret sets.
type Merger struct {
	policy ConflictPolicy
}

// New creates a new Merger with the given conflict policy.
func New(policy ConflictPolicy) (*Merger, error) {
	switch policy {
	case PolicySkip, PolicyOverwrite, PolicyError:
		return &Merger{policy: policy}, nil
	}
	return nil, fmt.Errorf("secretmergetag: unknown policy %q", policy)
}

// Merge combines multiple tag maps into one, applying the conflict policy.
// Each input is a map[envKey]map[tagKey]tagValue.
func (m *Merger) Merge(sources ...map[string]map[string]string) (map[string]map[string]string, error) {
	if len(sources) == 0 {
		return nil, errors.New("secretmergetag: no sources provided")
	}
	out := make(map[string]map[string]string)
	for _, src := range sources {
		for envKey, tags := range src {
			if _, exists := out[envKey]; !exists {
				out[envKey] = make(map[string]string)
			}
			for tk, tv := range tags {
				if _, conflict := out[envKey][tk]; conflict {
					switch m.policy {
					case PolicySkip:
						continue
					case PolicyOverwrite:
						out[envKey][tk] = tv
					case PolicyError:
						return nil, fmt.Errorf("secretmergetag: conflict on key %q tag %q", envKey, tk)
					}
				} else {
					out[envKey][tk] = tv
				}
			}
		}
	}
	return out, nil
}
