// Package secretreplace provides find-and-replace transformation over secret values.
package secretreplace

import (
	"errors"
	"strings"
)

// Rule defines a single find-and-replace operation applied to secret values.
type Rule struct {
	Find    string
	Replace string
}

// Replacer applies find-and-replace rules to secret values.
type Replacer struct {
	rules []Rule
}

// New creates a Replacer with the given rules.
func New(rules []Rule) (*Replacer, error) {
	for _, r := range rules {
		if r.Find == "" {
			return nil, errors.New("secretreplace: find string must not be empty")
		}
	}
	return &Replacer{rules: rules}, nil
}

// ParseRules parses rules from strings of the form "find=replace".
func ParseRules(raw []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(raw))
	for _, s := range raw {
		idx := strings.Index(s, "=")
		if idx < 1 {
			return nil, errors.New("secretreplace: invalid rule " + s + ": expected find=replace")
		}
		rules = append(rules, Rule{
			Find:    s[:idx],
			Replace: s[idx+1:],
		})
	}
	return rules, nil
}

// Apply returns a new map with all rules applied to every secret value.
func (r *Replacer) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretreplace: secrets must not be empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		for _, rule := range r.rules {
			v = strings.ReplaceAll(v, rule.Find, rule.Replace)
		}
		out[k] = v
	}
	return out, nil
}
