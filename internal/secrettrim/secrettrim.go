// Package secrettrim provides trimming of secret values by length boundaries.
package secrettrim

import (
	"errors"
	"fmt"
)

// ErrNoRules is returned when no trim rules are provided.
var ErrNoRules = errors.New("secrettrim: no rules provided")

// Rule defines a trim operation for a specific secret key.
type Rule struct {
	Key   string
	Start int
	End   int // -1 means until end of string
}

// Trimmer applies substring trim rules to secrets.
type Trimmer struct {
	rules []Rule
}

// New creates a Trimmer with the given rules.
func New(rules []Rule) (*Trimmer, error) {
	if len(rules) == 0 {
		return nil, ErrNoRules
	}
	for _, r := range rules {
		if r.Key == "" {
			return nil, errors.New("secrettrim: rule key must not be empty")
		}
		if r.Start < 0 {
			return nil, fmt.Errorf("secrettrim: rule %q has negative start", r.Key)
		}
		if r.End != -1 && r.End < r.Start {
			return nil, fmt.Errorf("secrettrim: rule %q has end before start", r.Key)
		}
	}
	return &Trimmer{rules: rules}, nil
}

// Apply returns a new map with trim rules applied to matching keys.
// Keys without a matching rule are passed through unchanged.
func (t *Trimmer) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secrettrim: secrets must not be empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, r := range t.rules {
		v, ok := out[r.Key]
		if !ok {
			continue
		}
		if r.Start >= len(v) {
			out[r.Key] = ""
			continue
		}
		end := r.End
		if end == -1 || end > len(v) {
			end = len(v)
		}
		out[r.Key] = v[r.Start:end]
	}
	return out, nil
}
