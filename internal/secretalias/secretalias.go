// Package secretalias provides aliasing support for secret keys,
// allowing vault secret keys to be referenced under alternate names.
package secretalias

import (
	"errors"
	"fmt"
	"strings"
)

// ErrNoAlias is returned when a lookup finds no alias for the given key.
var ErrNoAlias = errors.New("secretalias: no alias found")

// ErrEmptyKey is returned when an empty key is provided.
var ErrEmptyKey = errors.New("secretalias: key must not be empty")

// Aliaser maps original secret keys to one or more alias names.
type Aliaser struct {
	aliases map[string][]string // original -> aliases
	reverse map[string]string   // alias -> original
}

// New returns a new Aliaser.
func New() *Aliaser {
	return &Aliaser{
		aliases: make(map[string][]string),
		reverse: make(map[string]string),
	}
}

// Add registers alias names for an original key.
func (a *Aliaser) Add(original string, aliases ...string) error {
	if original == "" {
		return ErrEmptyKey
	}
	for _, alias := range aliases {
		if alias == "" {
			return ErrEmptyKey
		}
		if existing, ok := a.reverse[alias]; ok && existing != original {
			return fmt.Errorf("secretalias: alias %q already registered for %q", alias, existing)
		}
		a.reverse[alias] = original
	}
	a.aliases[original] = append(a.aliases[original], aliases...)
	return nil
}

// Aliases returns all aliases registered for the given original key.
func (a *Aliaser) Aliases(original string) ([]string, error) {
	if original == "" {
		return nil, ErrEmptyKey
	}
	v, ok := a.aliases[original]
	if !ok {
		return nil, ErrNoAlias
	}
	out := make([]string, len(v))
	copy(out, v)
	return out, nil
}

// Resolve returns the original key for a given alias.
func (a *Aliaser) Resolve(alias string) (string, error) {
	if alias == "" {
		return "", ErrEmptyKey
	}
	orig, ok := a.reverse[alias]
	if !ok {
		return "", ErrNoAlias
	}
	return orig, nil
}

// Apply expands secrets map with alias keys pointing to the same values.
func (a *Aliaser) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretalias: secrets must not be empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
		if aliases, ok := a.aliases[strings.ToUpper(k)]; ok {
			for _, alias := range aliases {
				out[alias] = v
			}
		}
		if aliases, ok := a.aliases[k]; ok {
			for _, alias := range aliases {
				out[alias] = v
			}
		}
	}
	return out, nil
}
