// Package diff provides utilities for comparing secret maps and
// reporting which keys were added, removed, or changed between runs.
package diff

import "sort"

// Result holds the outcome of comparing two secret snapshots.
type Result struct {
	Added   []string
	Removed []string
	Changed []string
}

// IsEmpty returns true when there are no differences.
func (r Result) IsEmpty() bool {
	return len(r.Added) == 0 && len(r.Removed) == 0 && len(r.Changed) == 0
}

// Compare returns a Result describing how current differs from previous.
// Both maps must use secret key strings as keys and plaintext values.
func Compare(previous, current map[string]string) Result {
	var res Result

	for k, cv := range current {
		pv, exists := previous[k]
		if !exists {
			res.Added = append(res.Added, k)
		} else if pv != cv {
			res.Changed = append(res.Changed, k)
		}
	}

	for k := range previous {
		if _, exists := current[k]; !exists {
			res.Removed = append(res.Removed, k)
		}
	}

	sort.Strings(res.Added)
	sort.Strings(res.Removed)
	sort.Strings(res.Changed)

	return res
}
