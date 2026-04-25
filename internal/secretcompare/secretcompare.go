// Package secretcompare provides side-by-side comparison of two secret maps,
// producing a structured result describing matches, mismatches, and missing keys.
package secretcompare

import (
	"errors"
	"fmt"	
	"sort"
)

// Status represents the comparison result for a single key.
type Status string

const (
	StatusMatch    Status = "match"
	StatusMismatch Status = "mismatch"
	StatusMissingA Status = "missing_a"
	StatusMissingB Status = "missing_b"
)

// Result holds the comparison outcome for one key.
type Result struct {
	Key    string
	Status Status
	ValueA string
	ValueB string
}

// Comparer compares two secret maps.
type Comparer struct{}

// New returns a new Comparer.
func New() *Comparer {
	return &Comparer{}
}

// Compare performs a key-by-key comparison of a and b.
// Returns an error if both maps are empty.
func (c *Comparer) Compare(a, b map[string]string) ([]Result, error) {
	if len(a) == 0 && len(b) == 0 {
		return nil, errors.New("secretcompare: both secret maps are empty")
	}

	keys := make(map[string]struct{})
	for k := range a {
		keys[k] = struct{}{}
	}
	for k := range b {
		keys[k] = struct{}{}
	}

	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	results := make([]Result, 0, len(sorted))
	for _, k := range sorted {
		va, inA := a[k]
		vb, inB := b[k]

		var status Status
		switch {
		case inA && inB && va == vb:
			status = StatusMatch
		case inA && inB:
			status = StatusMismatch
		case inA:
			status = StatusMissingB
		default:
			status = StatusMissingA
		}

		results = append(results, Result{
			Key:    k,
			Status: status,
			ValueA: va,
			ValueB: vb,
		})
	}
	return results, nil
}

// Summary returns a human-readable summary line for the comparison results.
func Summary(results []Result) string {
	var match, mismatch, missingA, missingB int
	for _, r := range results {
		switch r.Status {
		case StatusMatch:
			match++
		case StatusMismatch:
			mismatch++
		case StatusMissingA:
			missingA++
		case StatusMissingB:
			missingB++
		}
	}
	return fmt.Sprintf("match=%d mismatch=%d missing_a=%d missing_b=%d", match, mismatch, missingA, missingB)
}
