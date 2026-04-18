// Package secretsort provides utilities for sorting secret maps by key or value.
package secretsort

import (
	"errors"
	"sort"
	"strings"
)

// Order defines the sort direction.
type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

// ErrEmptySecrets is returned when the secrets map is nil or empty.
var ErrEmptySecrets = errors.New("secretsort: secrets map is empty")

// ErrInvalidOrder is returned when an unrecognised order string is provided.
var ErrInvalidOrder = errors.New("secretsort: invalid order, must be 'asc' or 'desc'")

// Sorter sorts a secrets map into an ordered slice of key-value pairs.
type Sorter struct {
	order Order
	byValue bool
}

// Entry is a single key-value pair from the sorted result.
type Entry struct {
	Key   string
	Value string
}

// New creates a Sorter. order must be "asc" or "desc". byValue controls
// whether sorting is performed on values instead of keys.
func New(order Order, byValue bool) (*Sorter, error) {
	o := Order(strings.ToLower(string(order)))
	if o != OrderAsc && o != OrderDesc {
		return nil, ErrInvalidOrder
	}
	return &Sorter{order: o, byValue: byValue}, nil
}

// Apply returns the secrets as a sorted slice of Entry.
func (s *Sorter) Apply(secrets map[string]string) ([]Entry, error) {
	if len(secrets) == 0 {
		return nil, ErrEmptySecrets
	}

	entries := make([]Entry, 0, len(secrets))
	for k, v := range secrets {
		entries = append(entries, Entry{Key: k, Value: v})
	}

	sort.Slice(entries, func(i, j int) bool {
		var a, b string
		if s.byValue {
			a, b = entries[i].Value, entries[j].Value
		} else {
			a, b = entries[i].Key, entries[j].Key
		}
		if s.order == OrderAsc {
			return a < b
		}
		return a > b
	})

	return entries, nil
}
