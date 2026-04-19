// Package secretmask provides pattern-based masking of secret values before output.
package secretmask

import (
	"errors"
	"regexp"
	"strings"
)

// Rule defines a pattern and replacement for masking.
type Rule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// Masker applies masking rules to secret values.
type Masker struct {
	rules []Rule
}

// New returns a Masker with the given rules.
func New(rules []Rule) *Masker {
	return &Masker{rules: rules}
}

// ParseRules parses rules from strings of the form "pattern=replacement".
func ParseRules(raw []string) ([]Rule, error) {
	var rules []Rule
	for _, r := range raw {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("secretmask: invalid rule format, expected pattern=replacement: " + r)
		}
		re, err := regexp.Compile(parts[0])
		if err != nil {
			return nil, errors.New("secretmask: invalid pattern '" + parts[0] + "': " + err.Error())
		}
		rules = append(rules, Rule{Pattern: re, Replacement: parts[1]})
	}
	return rules, nil
}

// Apply masks values in secrets according to the configured rules.
// Returns an error if secrets is empty.
func (m *Masker) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretmask: secrets must not be empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		masked := v
		for _, rule := range m.rules {
			masked = rule.Pattern.ReplaceAllString(masked, rule.Replacement)
		}
		out[k] = masked
	}
	return out, nil
}
