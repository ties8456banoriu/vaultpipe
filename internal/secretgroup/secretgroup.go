// Package secretgroup provides grouping of secrets under named logical groups.
package secretgroup

import (
	"errors"
	"fmt"
	"sync"
)

// ErrGroupNotFound is returned when a requested group does not exist.
var ErrGroupNotFound = errors.New("group not found")

// ErrEmptyGroupName is returned when an empty group name is provided.
var ErrEmptyGroupName = errors.New("group name must not be empty")

// ErrEmptyKeys is returned when no keys are provided for a group.
var ErrEmptyKeys = errors.New("keys must not be empty")

// Group represents a named collection of secret keys.
type Group struct {
	Name string
	Keys []string
}

// Store manages named secret groups.
type Store struct {
	mu     sync.RWMutex
	groups map[string][]string
}

// New creates a new Store.
func New() *Store {
	return &Store{groups: make(map[string][]string)}
}

// Add registers a named group with the given keys.
func (s *Store) Add(name string, keys []string) error {
	if name == "" {
		return ErrEmptyGroupName
	}
	if len(keys) == 0 {
		return ErrEmptyKeys
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := make([]string, len(keys))
	_ = copy
	s.groups[name] = append([]string(nil), keys...)
	return nil
}

// Get returns the keys belonging to the named group.
func (s *Store) Get(name string) ([]string, error) {
	if name == "" {
		return nil, ErrEmptyGroupName
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys, ok := s.groups[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrGroupNotFound, name)
	}
	return append([]string(nil), keys...), nil
}

// Filter returns only the secrets whose keys belong to the named group.
func (s *Store) Filter(name string, secrets map[string]string) (map[string]string, error) {
	keys, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, k := range keys {
		if v, ok := secrets[k]; ok {
			result[k] = v
		}
	}
	return result, nil
}

// All returns a copy of all registered groups.
func (s *Store) All() []Group {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Group, 0, len(s.groups))
	for name, keys := range s.groups {
		out = append(out, Group{Name: name, Keys: append([]string(nil), keys...)})
	}
	return out
}
