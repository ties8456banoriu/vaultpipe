// Package secretcopy provides functionality for copying secrets between
// namespaces or paths with optional key remapping.
package secretcopy

import (
	"errors"
	"fmt"
	"strings"
)

// Rule defines a copy operation from one key to another.
type Rule struct {
	From string
	To   string
}

// Copier copies secrets according to configured rules.
type Copier struct {
	rules []Rule
}

// New returns a new Copier with the given rules.
func New(rules []Rule) (*Copier, error) {
	if len(rules) == 0 {
		return nil, errors.New("secretcopy: at least one rule is required")
	}
	for _, r := range rules {
		if r.From == "" {
			return nil, errors.New("secretcopy: rule From key must not be empty")
		}
		if r.To == "" {
			return nil, errors.New("secretcopy: rule To key must not be empty")
		}
	}
	return &Copier{rules: rules}, nil
}

// ParseRules parses rules from strings in the format "FROM:TO".
func ParseRules(raw []string) ([]Rule, error) {
	var rules []Rule
	for _, s := range raw {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("secretcopy: invalid rule %q, expected FROM:TO", s)
		}
		rules = append(rules, Rule{From: parts[0], To: parts[1]})
	}
	return rules, nil
}

// Apply copies values in secrets according to rules, returning a new map.
// Source secrets are not modified. Missing From keys return an error.
func (c *Copier) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretcopy: secrets must not be empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, r := range c.rules {
		val, ok := secrets[r.From]
		if !ok {
			return nil, fmt.Errorf("secretcopy: source key %q not found", r.From)
		}
		out[r.To] = val
	}
	return out, nil
}
