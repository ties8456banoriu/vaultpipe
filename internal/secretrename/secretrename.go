// Package secretrename provides functionality for renaming secret keys
// in a secrets map using explicit rename rules.
package secretrename

import (
	"errors"
	"fmt"
	"strings"
)

// Rule represents a single rename rule mapping an old key to a new key.
type Rule struct {
	From string
	To   string
}

// Renamer applies rename rules to a secrets map.
type Renamer struct {
	rules []Rule
}

// New creates a new Renamer with the given rules.
func New(rules []Rule) (*Renamer, error) {
	if len(rules) == 0 {
		return nil, errors.New("secretrename: at least one rule is required")
	}
	for _, r := range rules {
		if r.From == "" {
			return nil, errors.New("secretrename: rule From key must not be empty")
		}
		if r.To == "" {
			return nil, errors.New("secretrename: rule To key must not be empty")
		}
	}
	return &Renamer{rules: rules}, nil
}

// ParseRules parses rename rules from a slice of "FROM:TO" strings.
func ParseRules(raw []string) ([]Rule, error) {
	var rules []Rule
	for _, s := range raw {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("secretrename: invalid rule %q: expected FROM:TO", s)
		}
		if parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("secretrename: invalid rule %q: FROM and TO must not be empty", s)
		}
		rules = append(rules, Rule{From: parts[0], To: parts[1]})
	}
	return rules, nil
}

// Apply renames keys in the secrets map according to the configured rules.
// Returns an error if secrets is empty or a From key is not found.
func (r *Renamer) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretrename: secrets must not be empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, rule := range r.rules {
		val, ok := out[rule.From]
		if !ok {
			return nil, fmt.Errorf("secretrename: key %q not found in secrets", rule.From)
		}
		delete(out, rule.From)
		out[rule.To] = val
	}
	return out, nil
}
