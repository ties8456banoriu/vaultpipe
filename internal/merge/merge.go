// Package merge provides utilities for combining secrets from multiple
// sources, with configurable precedence rules.
package merge

import "errors"

// Strategy defines how conflicts between sources are resolved.
type Strategy int

const (
	// StrategyFirst keeps the value from the first source that defines a key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last source that defines a key.
	StrategyLast
)

// ErrNoSources is returned when Merge is called with no sources.
var ErrNoSources = errors.New("merge: at least one source is required")

// ErrEmptySources is returned when all provided sources are nil or empty.
var ErrEmptySources = errors.New("merge: all sources are empty")

// Merger combines multiple secret maps according to a strategy.
type Merger struct {
	strategy Strategy
}

// NewMerger returns a Merger configured with the given strategy.
func NewMerger(strategy Strategy) *Merger {
	return &Merger{strategy: strategy}
}

// Merge combines the provided secret maps into a single map.
// Sources are processed in order; conflict resolution depends on the strategy.
// Returns an error if no sources are provided or all are empty.
func (m *Merger) Merge(sources ...map[string]string) (map[string]string, error) {
	if len(sources) == 0 {
		return nil, ErrNoSources
	}

	totalKeys := 0
	for _, s := range sources {
		totalKeys += len(s)
	}
	if totalKeys == 0 {
		return nil, ErrEmptySources
	}

	result := make(map[string]string, totalKeys)

	switch m.strategy {
	case StrategyLast:
		for _, src := range sources {
			for k, v := range src {
				result[k] = v
			}
		}
	default: // StrategyFirst
		for _, src := range sources {
			for k, v := range src {
				if _, exists := result[k]; !exists {
					result[k] = v
				}
			}
		}
	}

	return result, nil
}
