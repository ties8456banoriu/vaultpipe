// Package filter provides include/exclude filtering of secret keys
// before they are written to .env files.
//
// Rules are expressed as simple glob patterns (supporting '*' wildcard).
// Patterns prefixed with '!' are treated as exclusions.
//
// Example usage:
//
//	rules, err := filter.ParseRules([]string{"DB_*", "!DB_PASSWORD"})
//	f := filter.NewFilter(rules)
//	filtered, err := f.Apply(secrets)
package filter
