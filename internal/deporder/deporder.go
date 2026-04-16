// Package deporder resolves a deterministic write order for secrets
// based on declared dependencies between environment keys.
package deporder

import (
	"errors"
	"fmt"
)

// ErrCycle is returned when a dependency cycle is detected.
var ErrCycle = errors.New("deporder: dependency cycle detected")

// ErrUnknownKey is returned when a dependency references an undeclared key.
var ErrUnknownKey = errors.New("deporder: unknown dependency key")

// Resolver sorts secret keys into a stable write order respecting dependencies.
type Resolver struct {
	deps map[string][]string
}

// New returns a new Resolver. deps maps each key to the keys it depends on.
func New(deps map[string][]string) *Resolver {
	copy := make(map[string][]string, len(deps))
	for k, v := range deps {
		copy[k] = v
	}
	return &Resolver{deps: copy}
}

// Resolve returns keys from secrets in dependency order.
// Keys not present in deps are appended last in sorted order.
func (r *Resolver) Resolve(secrets map[string]string) ([]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("deporder: secrets must not be empty")
	}

	// Validate deps reference known keys.
	for key, deps := range r.deps {
		_ = key
		for _, d := range deps {
			if _, ok := secrets[d]; !ok {
				if _, ok2 := r.deps[d]; !ok2 {
					return nil, fmt.Errorf("%w: %q", ErrUnknownKey, d)
				}
			}
		}
	}

	visited := make(map[string]bool)
	inStack := make(map[string]bool)
	var order []string

	var visit func(k string) error
	visit = func(k string) error {
		if inStack[k] {
			return fmt.Errorf("%w: key %q", ErrCycle, k)
		}
		if visited[k] {
			return nil
		}
		inStack[k] = true
		for _, dep := range r.deps[k] {
			if err := visit(dep); err != nil {
				return err
			}
		}
		inStack[k] = false
		visited[k] = true
		order = append(order, k)
		return nil
	}

	for k := range secrets {
		if err := visit(k); err != nil {
			return nil, err
		}
	}
	return order, nil
}
