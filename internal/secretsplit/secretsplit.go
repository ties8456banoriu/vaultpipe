// Package secretsplit provides functionality for splitting a secrets map
// into multiple named subsets based on key prefix rules.
package secretsplit

import (
	"errors"
	"strings"
)

// ErrEmptySecrets is returned when the input secrets map is empty.
var ErrEmptySecrets = errors.New("secretsplit: secrets map is empty")

// ErrNoRules is returned when no split rules are provided.
var ErrNoRules = errors.New("secretsplit: no split rules provided")

// ErrEmptyName is returned when a rule has an empty subset name.
var ErrEmptyName = errors.New("secretsplit: rule has empty subset name")

// Rule defines a named subset and the key prefix used to select secrets.
type Rule struct {
	Name   string
	Prefix string
}

// Splitter splits a secrets map into named subsets.
type Splitter struct{}

// New returns a new Splitter.
func New() *Splitter {
	return &Splitter{}
}

// Apply splits secrets into named subsets according to the given rules.
// Keys are matched by prefix (case-insensitive). A key may appear in multiple
// subsets if multiple rules match. Keys not matched by any rule are placed
// under the reserved name "_unmatched".
func (s *Splitter) Apply(secrets map[string]string, rules []Rule) (map[string]map[string]string, error) {
	if len(secrets) == 0 {
		return nil, ErrEmptySecrets
	}
	if len(rules) == 0 {
		return nil, ErrNoRules
	}
	for _, r := range rules {
		if r.Name == "" {
			return nil, ErrEmptyName
		}
	}

	result := make(map[string]map[string]string)
	matched := make(map[string]bool)

	for _, rule := range rules {
		subset := make(map[string]string)
		prefix := strings.ToUpper(rule.Prefix)
		for k, v := range secrets {
			if strings.HasPrefix(strings.ToUpper(k), prefix) {
				subset[k] = v
				matched[k] = true
			}
		}
		result[rule.Name] = subset
	}

	unmatched := make(map[string]string)
	for k, v := range secrets {
		if !matched[k] {
			unmatched[k] = v
		}
	}
	if len(unmatched) > 0 {
		result["_unmatched"] = unmatched
	}

	return result, nil
}
