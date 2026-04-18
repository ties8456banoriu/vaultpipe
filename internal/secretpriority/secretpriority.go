// Package secretpriority assigns and enforces priority levels to secrets,
// allowing higher-priority sources to override lower-priority ones.
package secretpriority

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// Level represents a numeric priority (higher = more important).
type Level int

// Entry holds a secret value alongside its priority.
type Entry struct {
	Value    string
	Priority Level
}

// Manager tracks priority entries per env key.
type Manager struct {
	mu      sync.RWMutex
	entries map[string][]Entry
}

// New returns an initialised Manager.
func New() *Manager {
	return &Manager{entries: make(map[string][]Entry)}
}

// Add records a value for key at the given priority level.
func (m *Manager) Add(key, value string, priority Level) error {
	if key == "" {
		return errors.New("secretpriority: key must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[key] = append(m.entries[key], Entry{Value: value, Priority: priority})
	return nil
}

// Resolve returns the value with the highest priority for key.
// If multiple entries share the top priority the first-added wins.
func (m *Manager) Resolve(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	entries, ok := m.entries[key]
	if !ok || len(entries) == 0 {
		return "", fmt.Errorf("secretpriority: no entries for key %q", key)
	}
	best := entries[0]
	for _, e := range entries[1:] {
		if e.Priority > best.Priority {
			best = e
		}
	}
	return best.Value, nil
}

// ResolveAll resolves every tracked key and returns the winning map.
func (m *Manager) ResolveAll() (map[string]string, error) {
	m.mu.RLock()
	keys := make([]string, 0, len(m.entries))
	for k := range m.entries {
		keys = append(keys, k)
	}
	m.mu.RUnlock()

	if len(keys) == 0 {
		return nil, errors.New("secretpriority: no secrets tracked")
	}
	sort.Strings(keys)
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		v, err := m.Resolve(k)
		if err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, nil
}
