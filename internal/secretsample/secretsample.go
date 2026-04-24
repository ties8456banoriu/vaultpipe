// Package secretsample provides random sampling of secrets from a map.
// It is useful for testing, auditing, or previewing a subset of secrets
// without exposing the full set.
package secretsample

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

// ErrEmptySecrets is returned when the input secrets map is empty.
var ErrEmptySecrets = errors.New("secretsample: secrets map is empty")

// ErrInvalidN is returned when the requested sample size is less than 1.
var ErrInvalidN = errors.New("secretsample: n must be at least 1")

// Sampler draws a random sample of secrets.
type Sampler struct {
	n    int
	seed int64
}

// Option is a functional option for Sampler.
type Option func(*Sampler)

// WithSeed sets a deterministic random seed for reproducible sampling.
func WithSeed(seed int64) Option {
	return func(s *Sampler) {
		s.seed = seed
	}
}

// New creates a Sampler that returns up to n secrets per call.
// Returns ErrInvalidN if n < 1.
func New(n int, opts ...Option) (*Sampler, error) {
	if n < 1 {
		return nil, fmt.Errorf("%w: got %d", ErrInvalidN, n)
	}
	s := &Sampler{n: n, seed: rand.Int63()}
	for _, o := range opts {
		o(s)
	}
	return s, nil
}

// Apply returns a random sample of up to n entries from secrets.
// The returned map is a new map; the original is not modified.
// Returns ErrEmptySecrets if secrets is empty.
func (s *Sampler) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, ErrEmptySecrets
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	// Sort for determinism before shuffling.
	sort.Strings(keys)

	r := rand.New(rand.NewSource(s.seed))
	r.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })

	pick := s.n
	if pick > len(keys) {
		pick = len(keys)
	}

	out := make(map[string]string, pick)
	for _, k := range keys[:pick] {
		out[k] = secrets[k]
	}
	return out, nil
}
