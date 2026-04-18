// Package secretclone provides functionality for cloning (deep-copying)
// a secrets map with optional key transformations.
package secretclone

import (
	"errors"
	"strings"
)

// ErrEmptySecrets is returned when an empty secrets map is provided.
var ErrEmptySecrets = errors.New("secretclone: secrets map is empty")

// Option configures the Cloner.
type Option func(*Cloner)

// WithPrefix prepends a prefix to every key in the cloned map.
func WithPrefix(prefix string) Option {
	return func(c *Cloner) { c.prefix = prefix }
}

// WithUppercase uppercases all keys in the cloned map.
func WithUppercase() Option {
	return func(c *Cloner) { c.uppercase = true }
}

// Cloner clones secrets maps.
type Cloner struct {
	prefix    string
	uppercase bool
}

// New returns a new Cloner configured with the given options.
func New(opts ...Option) *Cloner {
	c := &Cloner{}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Clone returns a deep copy of src, applying any configured transformations.
func (c *Cloner) Clone(src map[string]string) (map[string]string, error) {
	if len(src) == 0 {
		return nil, ErrEmptySecrets
	}
	out := make(map[string]string, len(src))
	for k, v := range src {
		if c.uppercase {
			k = strings.ToUpper(k)
		}
		if c.prefix != "" {
			k = c.prefix + k
		}
		out[k] = v
	}
	return out, nil
}
