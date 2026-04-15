// Package validate provides pre-write validation of secret key-value pairs
// before they are injected into .env files.
package validate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ErrEmptySecrets is returned when the secrets map is nil or empty.
var ErrEmptySecrets = errors.New("validate: secrets map is empty")

// validKeyPattern matches valid environment variable names.
var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Result holds the outcome of validating a single key-value pair.
type Result struct {
	Key     string
	Warning string
}

// Validator checks secrets for common issues before writing.
type Validator struct {
	warnOnEmpty bool
}

// NewValidator returns a Validator. When warnOnEmpty is true, keys with
// empty string values produce a warning rather than an error.
func NewValidator(warnOnEmpty bool) *Validator {
	return &Validator{warnOnEmpty: warnOnEmpty}
}

// Validate checks all key-value pairs in secrets and returns any warnings
// and a combined error for hard failures.
func (v *Validator) Validate(secrets map[string]string) ([]Result, error) {
	if len(secrets) == 0 {
		return nil, ErrEmptySecrets
	}

	var warnings []Result
	var errs []string

	for k, val := range secrets {
		if !validKeyPattern.MatchString(k) {
			errs = append(errs, fmt.Sprintf("invalid key name %q: must match [A-Za-z_][A-Za-z0-9_]*", k))
			continue
		}
		if val == "" && v.warnOnEmpty {
			warnings = append(warnings, Result{Key: k, Warning: "value is empty"})
		}
	}

	if len(errs) > 0 {
		return warnings, fmt.Errorf("validate: %s", strings.Join(errs, "; "))
	}
	return warnings, nil
}
