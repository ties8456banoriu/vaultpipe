// Package secretflatten provides a way to flatten nested secret keys
// by replacing a delimiter with a target separator, producing flat env-style keys.
package secretflatten

import (
	"errors"
	"strings"
)

// Flattener flattens secret keys by replacing a delimiter with a separator.
type Flattener struct {
	delimiter string
	separator string
	uppercase bool
}

// Option configures a Flattener.
type Option func(*Flattener)

// WithUppercase causes all output keys to be uppercased.
func WithUppercase() Option {
	return func(f *Flattener) { f.uppercase = true }
}

// New creates a Flattener that replaces delimiter with separator in all keys.
// delimiter and separator must be non-empty.
func New(delimiter, separator string, opts ...Option) (*Flattener, error) {
	if delimiter == "" {
		return nil, errors.New("secretflatten: delimiter must not be empty")
	}
	if separator == "" {
		return nil, errors.New("secretflatten: separator must not be empty")
	}
	f := &Flattener{delimiter: delimiter, separator: separator}
	for _, o := range opts {
		o(f)
	}
	return f, nil
}

// Apply returns a new map with all keys flattened.
// Returns an error if secrets is empty.
func (f *Flattener) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretflatten: secrets must not be empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		newKey := strings.ReplaceAll(k, f.delimiter, f.separator)
		if f.uppercase {
			newKey = strings.ToUpper(newKey)
		}
		out[newKey] = v
	}
	return out, nil
}
