// Package transform provides value transformation rules for secrets
// before they are written to .env files.
package transform

import (
	"fmt"
	"strings"
)

// Rule represents a single transformation to apply to a secret value.
type Rule struct {
	Key  string
	Type string // "upper", "lower", "trim", "prefix", "suffix"
	Arg  string // optional argument (e.g. prefix/suffix string)
}

// Transformer applies transformation rules to a map of secrets.
type Transformer struct {
	rules []Rule
}

// NewTransformer creates a Transformer with the given rules.
func NewTransformer(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply applies all matching transformation rules to the secrets map.
// Returns a new map with transformed values.
func (t *Transformer) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, fmt.Errorf("transform: secrets map is empty")
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}

	for _, rule := range t.rules {
		val, ok := result[rule.Key]
		if !ok {
			continue
		}
		transformed, err := applyRule(val, rule)
		if err != nil {
			return nil, fmt.Errorf("transform: key %q: %w", rule.Key, err)
		}
		result[rule.Key] = transformed
	}

	return result, nil
}

// ParseRules parses a slice of raw rule strings in the format "KEY:TYPE[:ARG]".
func ParseRules(raw []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(raw))
	for _, s := range raw {
		parts := strings.SplitN(s, ":", 3)
		if len(parts) < 2 {
			return nil, fmt.Errorf("transform: invalid rule %q, expected KEY:TYPE[:ARG]", s)
		}
		rule := Rule{Key: parts[0], Type: parts[1]}
		if len(parts) == 3 {
			rule.Arg = parts[2]
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func applyRule(val string, rule Rule) (string, error) {
	switch rule.Type {
	case "upper":
		return strings.ToUpper(val), nil
	case "lower":
		return strings.ToLower(val), nil
	case "trim":
		return strings.TrimSpace(val), nil
	case "prefix":
		return rule.Arg + val, nil
	case "suffix":
		return val + rule.Arg, nil
	default:
		return "", fmt.Errorf("unknown transform type %q", rule.Type)
	}
}
