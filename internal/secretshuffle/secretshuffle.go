// Package secretshuffle provides randomised ordering of secret maps.
package secretshuffle

import (
	"errors"
	"math/rand"
)

// Shuffler randomises the iteration order of a secrets map by returning
// an ordered slice of keys in a random sequence.
type Shuffler struct {
	rng *rand.Rand
}

// New returns a Shuffler seeded with the provided source.
// Pass nil to use a default rand source.
func New(src rand.Source) *Shuffler {
	if src == nil {
		src = rand.NewSource(rand.Int63())
	}
	return &Shuffler{rng: rand.New(src)}
}

// Shuffle returns the keys of secrets in a randomised order.
// Returns an error if secrets is empty.
func (s *Shuffler) Shuffle(secrets map[string]string) ([]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretshuffle: secrets must not be empty")
	}
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	s.rng.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	return keys, nil
}

// Apply returns a new map built by copying secrets using the shuffled key
// order. Because map iteration is unordered in Go the returned map is
// identical in content; the useful artefact is the accompanying key slice.
func (s *Shuffler) Apply(secrets map[string]string) (map[string]string, []string, error) {
	keys, err := s.Shuffle(secrets)
	if err != nil {
		return nil, nil, err
	}
	out := make(map[string]string, len(secrets))
	for _, k := range keys {
		out[k] = secrets[k]
	}
	return out, keys, nil
}
