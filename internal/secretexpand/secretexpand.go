// Package secretexpand provides variable interpolation within secret values,
// allowing secrets to reference other secrets using ${KEY} syntax.
package secretexpand

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrEmptySecrets = errors.New("secretexpand: secrets map is empty")
var ErrCircularReference = errors.New("secretexpand: circular reference detected")

var refPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

// Expander interpolates ${KEY} references within secret values.
type Expander struct{}

// New returns a new Expander.
func New() *Expander {
	return &Expander{}
}

// Apply resolves all ${KEY} references in the secrets map.
// Returns an error if the map is empty, a key is unknown, or a circular
// reference is detected.
func (e *Expander) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, ErrEmptySecrets
	}

	resolved := make(map[string]string, len(secrets))
	for k, v := range secrets {
		resolved[k] = v
	}

	for key := range resolved {
		if err := resolve(resolved, key, []string{}); err != nil {
			return nil, err
		}
	}
	return resolved, nil
}

func resolve(secrets map[string]string, key string, stack []string) error {
	for _, s := range stack {
		if s == key {
			return fmt.Errorf("%w: %s", ErrCircularReference, strings.Join(append(stack, key), " -> "))
		}
	}

	val := secrets[key]
	matches := refPattern.FindAllStringSubmatch(val, -1)
	if len(matches) == 0 {
		return nil
	}

	for _, m := range matches {
		refKey := m[1]
		refVal, ok := secrets[refKey]
		if !ok {
			return fmt.Errorf("secretexpand: unknown reference ${%s} in key %s", refKey, key)
		}
		if err := resolve(secrets, refKey, append(stack, key)); err != nil {
			return err
		}
		val = strings.ReplaceAll(val, m[0], refVal)
	}
	secrets[key] = val
	return nil
}
