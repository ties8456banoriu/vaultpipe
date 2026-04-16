// Package truncate provides utilities for truncating secret values
// to a maximum byte length before writing to .env files.
package truncate

import "errors"

// ErrEmptySecrets is returned when an empty secrets map is provided.
var ErrEmptySecrets = errors.New("truncate: secrets map is empty")

// ErrInvalidMaxLen is returned when maxLen is less than 1.
var ErrInvalidMaxLen = errors.New("truncate: maxLen must be at least 1")

// Truncator applies a maximum byte length to secret values.
type Truncator struct {
	maxLen int
}

// New creates a Truncator that will truncate values to at most maxLen bytes.
func New(maxLen int) (*Truncator, error) {
	if maxLen < 1 {
		return nil, ErrInvalidMaxLen
	}
	return &Truncator{maxLen: maxLen}, nil
}

// Apply returns a new map with all values truncated to the configured length.
// The original map is not modified.
func (t *Truncator) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, ErrEmptySecrets
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if len(v) > t.maxLen {
			out[k] = v[:t.maxLen]
		} else {
			out[k] = v
		}
	}
	return out, nil
}
