// Package filter provides functionality for selectively including or
// excluding secrets based on key patterns before they are written to .env files.
package filter

import (
	"fmt"
	"strings"
)

// Rule defines a single include or exclude pattern for secret keys.
type Rule struct {
	Pattern string
	Exclude bool
}

// Filter applies include/exclude rules to a map of secrets.
type Filter struct {
	rules []Rule
}

// NewFilter creates a Filter from a slice of rules.
func NewFilter(rules []Rule) *Filter {
	return &Filter{rules: rules}
}

// ParseRules parses a slice of rule strings into Rule values.
// Strings prefixed with '!' are treated as exclude rules.
// Example: ["DB_*", "!DB_PASSWORD"] includes all DB_ keys except DB_PASSWORD.
func ParseRules(raw []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(raw))
	for _, r := range raw {
		r = strings.TrimSpace(r)
		if r == "" {
			return nil, fmt.Errorf("filter: empty rule is not allowed")
		}
		if strings.HasPrefix(r, "!") {
			rules = append(rules, Rule{Pattern: r[1:], Exclude: true})
		} else {
			rules = append(rules, Rule{Pattern: r, Exclude: false})
		}
	}
	return rules, nil
}

// Apply returns a filtered copy of secrets. If no include rules are present,
// all keys are included by default. Exclude rules are always applied last.
func (f *Filter) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, fmt.Errorf("filter: secrets map is empty")
	}

	hasInclude := false
	for _, r := range f.rules {
		if !r.Exclude {
			hasInclude = true
			break
		}
	}

	result := make(map[string]string)

	for k, v := range secrets {
		included := !hasInclude // if no include rules, include everything
		for _, r := range f.rules {
			if !r.Exclude && matchPattern(r.Pattern, k) {
				included = true
			}
		}
		if included {
			result[k] = v
		}
	}

	// Apply exclude rules
	for _, r := range f.rules {
		if r.Exclude {
			for k := range result {
				if matchPattern(r.Pattern, k) {
					delete(result, k)
				}
			}
		}
	}

	return result, nil
}

// matchPattern checks whether key matches a simple glob pattern (only '*' wildcard supported).
func matchPattern(pattern, key string) bool {
	if pattern == "*" {
		return true
	}
	if !strings.Contains(pattern, "*") {
		return pattern == key
	}
	parts := strings.SplitN(pattern, "*", 2)
	return strings.HasPrefix(key, parts[0]) && strings.HasSuffix(key, parts[1])
}
