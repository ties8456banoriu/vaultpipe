// Package secretpromote provides functionality for promoting secrets
// from one environment namespace to another (e.g., staging → production).
package secretpromote

import (
	"errors"
	"fmt"
	"strings"
)

// ErrEmptySource is returned when the source namespace is empty.
var ErrEmptySource = errors.New("secretpromote: source namespace must not be empty")

// ErrEmptyTarget is returned when the target namespace is empty.
var ErrEmptyTarget = errors.New("secretpromote: target namespace must not be empty")

// ErrEmptySecrets is returned when the input secrets map is empty.
var ErrEmptySecrets = errors.New("secretpromote: secrets must not be empty")

// ErrSameNamespace is returned when source and target are identical.
var ErrSameNamespace = errors.New("secretpromote: source and target namespace must differ")

// Result holds the outcome of a promotion operation.
type Result struct {
	Promoted map[string]string
	Skipped  []string
}

// Promoter copies secrets from one namespace prefix to another.
type Promoter struct {
	source string
	target string
	overwrite bool
}

// WithOverwrite configures the Promoter to overwrite existing target keys.
func WithOverwrite() func(*Promoter) {
	return func(p *Promoter) {
		p.overwrite = true
	}
}

// New creates a new Promoter that promotes secrets from source to target namespace.
func New(source, target string, opts ...func(*Promoter)) (*Promoter, error) {
	source = strings.TrimSpace(source)
	target = strings.TrimSpace(target)
	if source == "" {
		return nil, ErrEmptySource
	}
	if target == "" {
		return nil, ErrEmptyTarget
	}
	if source == target {
		return nil, ErrSameNamespace
	}
	p := &Promoter{source: source, target: target}
	for _, o := range opts {
		o(p)
	}
	return p, nil
}

// Apply promotes matching secrets from source namespace to target namespace.
// Keys prefixed with source are rewritten to use the target prefix.
func (p *Promoter) Apply(secrets map[string]string) (*Result, error) {
	if len(secrets) == 0 {
		return nil, ErrEmptySecrets
	}

	result := &Result{
		Promoted: make(map[string]string),
	}

	prefix := p.source + "_"
	for k, v := range secrets {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		suffix := strings.TrimPrefix(k, prefix)
		newKey := fmt.Sprintf("%s_%s", p.target, suffix)
		if _, exists := secrets[newKey]; exists && !p.overwrite {
			result.Skipped = append(result.Skipped, newKey)
			continue
		}
		result.Promoted[newKey] = v
	}

	return result, nil
}
