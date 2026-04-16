// Package schema provides secret key validation against a declared schema.
package schema

import (
	"errors"
	"fmt"
	"strings"
)

// FieldRule declares expectations for a single secret key.
type FieldRule struct {
	Key      string
	Required bool
	Pattern  string // optional: "string", "int", "bool"
}

// Validator checks secrets against a set of field rules.
type Validator struct {
	rules map[string]FieldRule
}

// NewValidator builds a Validator from a slice of FieldRules.
func NewValidator(rules []FieldRule) *Validator {
	m := make(map[string]FieldRule, len(rules))
	for _, r := range rules {
		m[r.Key] = r
	}
	return &Validator{rules: m}
}

// ParseRules parses rules from strings like "KEY:required:int".
func ParseRules(raw []string) ([]FieldRule, error) {
	var rules []FieldRule
	for _, s := range raw {
		parts := strings.Split(s, ":")
		if len(parts) < 1 || parts[0] == "" {
			return nil, fmt.Errorf("schema: invalid rule %q", s)
		}
		r := FieldRule{Key: parts[0]}
		for _, mod := range parts[1:] {
			switch mod {
			case "required":
				r.Required = true
			case "string", "int", "bool":
				r.Pattern = mod
			default:
				return nil, fmt.Errorf("schema: unknown modifier %q in rule %q", mod, s)
			}
		}
		rules = append(rules, r)
	}
	return rules, nil
}

// Validate checks secrets against declared rules. Returns all violations.
func (v *Validator) Validate(secrets map[string]string) error {
	var errs []string
	for key, rule := range v.rules {
		val, ok := secrets[key]
		if rule.Required && (!ok || val == "") {
			errs = append(errs, fmt.Sprintf("key %q is required but missing or empty", key))
			continue
		}
		if ok && rule.Pattern != "" {
			if err := checkPattern(val, rule.Pattern); err != nil {
				errs = append(errs, fmt.Sprintf("key %q: %s", key, err))
			}
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func checkPattern(val, pattern string) error {
	switch pattern {
	case "int":
		for _, c := range val {
			if c < '0' || c > '9' {
				return fmt.Errorf("expected int, got %q", val)
			}
		}
	case "bool":
		v := strings.ToLower(val)
		if v != "true" && v != "false" {
			return fmt.Errorf("expected bool (true/false), got %q", val)
		}
	}
	return nil
}
