// Package secretomit provides a filter that removes secrets whose values
// match a set of patterns, such as empty strings or placeholder values.
package secretomit

import (
	"errors"
	"regexp"
)

// Omitter removes secrets whose values match any registered pattern.
type Omitter struct {
	patterns []*regexp.Regexp
}

// New creates an Omitter from a slice of regex pattern strings.
// Returns an error if any pattern is invalid or the slice is empty.
func New(patterns []string) (*Omitter, error) {
	if len(patterns) == 0 {
		return nil, errors.New("secretomit: at least one pattern is required")
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if p == "" {
			return nil, errors.New("secretomit: pattern must not be empty")
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("secretomit: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &Omitter{patterns: compiled}, nil
}

// Apply returns a new map with secrets removed whose values match any pattern.
// Returns an error if secrets is empty.
func (o *Omitter) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretomit: secrets must not be empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if !o.matchesAny(v) {
			out[k] = v
		}
	}
	return out, nil
}

func (o *Omitter) matchesAny(value string) bool {
	for _, re := range o.patterns {
		if re.MatchString(value) {
			return true
		}
	}
	return false
}
