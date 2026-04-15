package mapping

import (
	"fmt"
	"strings"
)

// Rule defines a mapping from a Vault secret key to an env variable name.
type Rule struct {
	VaultKey string
	EnvKey   string
}

// Mapper holds a set of mapping rules for transforming Vault secret keys
// into environment variable names.
type Mapper struct {
	rules map[string]string
}

// NewMapper creates a Mapper from a slice of Rules.
func NewMapper(rules []Rule) *Mapper {
	m := &Mapper{rules: make(map[string]string, len(rules))}
	for _, r := range rules {
		m.rules[r.VaultKey] = r.EnvKey
	}
	return m
}

// Apply transforms a map of Vault secret key/values into env variable key/values
// using the configured rules. Keys without a mapping rule are converted to
// upper-snake-case by default.
func (m *Mapper) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, fmt.Errorf("mapping: secrets map is empty")
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		envKey, ok := m.rules[k]
		if !ok {
			envKey = defaultEnvKey(k)
		}
		result[envKey] = v
	}
	return result, nil
}

// ParseRules parses a slice of "vault_key=ENV_KEY" strings into Rules.
func ParseRules(raw []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(raw))
	for _, entry := range raw {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("mapping: invalid rule %q, expected format vault_key=ENV_KEY", entry)
		}
		rules = append(rules, Rule{VaultKey: parts[0], EnvKey: parts[1]})
	}
	return rules, nil
}

// Lookup returns the env key mapped to the given Vault key, falling back to
// the default upper-snake-case conversion if no explicit rule exists.
func (m *Mapper) Lookup(vaultKey string) string {
	if envKey, ok := m.rules[vaultKey]; ok {
		return envKey
	}
	return defaultEnvKey(vaultKey)
}

// defaultEnvKey converts a vault key to an upper-snake-case env variable name.
func defaultEnvKey(key string) string {
	replacer := strings.NewReplacer("-", "_", ".", "_", "/", "_")
	return strings.ToUpper(replacer.Replace(key))
}
