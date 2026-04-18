// Package secretsanitize provides value sanitization for secrets before
// they are written to .env files.
package secretsanitize

import (
	"errors"
	"strings"
)

// Rule defines a sanitization operation to apply to secret values.
type Rule struct {
	StripPrefix string
	StripSuffix string
	TrimSpace   bool
	ReplaceNewlines bool
}

// Sanitizer applies sanitization rules to a map of secrets.
type Sanitizer struct {
	rule Rule
}

// New returns a Sanitizer configured with the given Rule.
func New(r Rule) *Sanitizer {
	return &Sanitizer{rule: r}
}

// Apply returns a new map with sanitized values. Returns an error if secrets is empty.
func (s *Sanitizer) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretsanitize: secrets map is empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = s.sanitize(v)
	}
	return out, nil
}

func (s *Sanitizer) sanitize(v string) string {
	if s.rule.TrimSpace {
		v = strings.TrimSpace(v)
	}
	if s.rule.StripPrefix != "" {
		v = strings.TrimPrefix(v, s.rule.StripPrefix)
	}
	if s.rule.StripSuffix != "" {
		v = strings.TrimSuffix(v, s.rule.StripSuffix)
	}
	if s.rule.ReplaceNewlines {
		v = strings.ReplaceAll(v, "\n", " ")
		v = strings.ReplaceAll(v, "\r", "")
	}
	return v
}
