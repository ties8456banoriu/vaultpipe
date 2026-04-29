// Package secretwipe provides functionality for zeroing out secret values
// in memory after use, reducing the window of exposure for sensitive data.
package secretwipe

import "errors"

// Wiper zeroes secret values in a map after use.
type Wiper struct {
	patterns []string
}

// Option configures a Wiper.
type Option func(*Wiper)

// WithPatterns restricts wiping to keys matching any of the given glob-style
// prefixes. If no patterns are provided, all keys are wiped.
func WithPatterns(patterns []string) Option {
	return func(w *Wiper) {
		w.patterns = patterns
	}
}

// New creates a new Wiper with the given options.
func New(opts ...Option) *Wiper {
	w := &Wiper{}
	for _, o := range opts {
		o(w)
	}
	return w
}

// Wipe zeroes the values of secrets in place. If patterns are configured,
// only keys matching at least one pattern are wiped. Returns an error if
// secrets is nil or empty.
func (w *Wiper) Wipe(secrets map[string]string) error {
	if len(secrets) == 0 {
		return errors.New("secretwipe: no secrets to wipe")
	}

	for k := range secrets {
		if w.shouldWipe(k) {
			secrets[k] = ""
		}
	}
	return nil
}

// WipeCopy returns a copy of secrets with matching values zeroed, leaving the
// original map untouched.
func (w *Wiper) WipeCopy(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretwipe: no secrets to wipe")
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if w.shouldWipe(k) {
			out[k] = ""
		} else {
			out[k] = v
		}
	}
	return out, nil
}

func (w *Wiper) shouldWipe(key string) bool {
	if len(w.patterns) == 0 {
		return true
	}
	for _, p := range w.patterns {
		if matchPrefix(key, p) {
			return true
		}
	}
	return false
}

func matchPrefix(key, pattern string) bool {
	if pattern == "*" {
		return true
	}
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		return len(key) >= len(pattern)-1 && key[:len(pattern)-1] == pattern[:len(pattern)-1]
	}
	return key == pattern
}
