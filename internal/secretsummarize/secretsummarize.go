// Package secretsummarize provides a summarizer that produces a concise
// statistical overview of a secret map: key count, empty values, average
// value length, and the longest/shortest keys.
package secretsummarize

import (
	"errors"
	"sort"
)

// Summary holds aggregate statistics about a set of secrets.
type Summary struct {
	TotalKeys    int
	EmptyValues  int
	AvgValueLen  float64
	LongestKey   string
	ShortestKey  string
	UniqueValues int
}

// Summarizer computes statistics over a secret map.
type Summarizer struct{}

// New returns a new Summarizer.
func New() *Summarizer {
	return &Summarizer{}
}

// Apply computes a Summary from the provided secrets map.
// It returns an error if secrets is nil or empty.
func (s *Summarizer) Apply(secrets map[string]string) (Summary, error) {
	if len(secrets) == 0 {
		return Summary{}, errors.New("secretsummarize: secrets must not be empty")
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var totalLen int
	empty := 0
	seen := make(map[string]struct{})

	for _, k := range keys {
		v := secrets[k]
		if v == "" {
			empty++
		}
		totalLen += len(v)
		seen[v] = struct{}{}
	}

	longest := keys[0]
	shortest := keys[0]
	for _, k := range keys[1:] {
		if len(k) > len(longest) {
			longest = k
		}
		if len(k) < len(shortest) {
			shortest = k
		}
	}

	avg := 0.0
	if len(keys) > 0 {
		avg = float64(totalLen) / float64(len(keys))
	}

	return Summary{
		TotalKeys:    len(keys),
		EmptyValues:  empty,
		AvgValueLen:  avg,
		LongestKey:   longest,
		ShortestKey:  shortest,
		UniqueValues: len(seen),
	}, nil
}
