// Package secretnamespace provides namespace-based isolation for secrets,
// allowing secrets to be grouped and accessed under named namespaces.
package secretnamespace

import (
	"errors"
	"fmt"
	"sync"
)

// ErrNamespaceNotFound is returned when a namespace does not exist.
var ErrNamespaceNotFound = errors.New("namespace not found")

// ErrEmptyNamespace is returned when an empty namespace name is provided.
var ErrEmptyNamespace = errors.New("namespace name must not be empty")

// ErrEmptySecrets is returned when an empty secrets map is provided.
var ErrEmptySecrets = errors.New("secrets must not be empty")

// Manager manages namespaced secrets.
type Manager struct {
	mu         sync.RWMutex
	namespaces map[string]map[string]string
}

// New creates a new namespace Manager.
func New() *Manager {
	return &Manager{
		namespaces: make(map[string]map[string]string),
	}
}

// Set stores secrets under the given namespace, replacing any existing entry.
func (m *Manager) Set(namespace string, secrets map[string]string) error {
	if namespace == "" {
		return ErrEmptyNamespace
	}
	if len(secrets) == 0 {
		return ErrEmptySecrets
	}
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.namespaces[namespace] = copy
	return nil
}

// Get returns the secrets stored under the given namespace.
func (m *Manager) Get(namespace string) (map[string]string, error) {
	if namespace == "" {
		return nil, ErrEmptyNamespace
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	secrets, ok := m.namespaces[namespace]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNamespaceNotFound, namespace)
	}
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	return copy, nil
}

// Delete removes a namespace and its secrets.
func (m *Manager) Delete(namespace string) error {
	if namespace == "" {
		return ErrEmptyNamespace
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.namespaces[namespace]; !ok {
		return fmt.Errorf("%w: %s", ErrNamespaceNotFound, namespace)
	}
	delete(m.namespaces, namespace)
	return nil
}

// List returns all registered namespace names.
func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.namespaces))
	for k := range m.namespaces {
		names = append(names, k)
	}
	return names
}

// Merge combines secrets from all namespaces into a single map.
// Later namespaces overwrite earlier ones on key conflict.
func (m *Manager) Merge(namespaces ...string) (map[string]string, error) {
	result := make(map[string]string)
	for _, ns := range namespaces {
		secrets, err := m.Get(ns)
		if err != nil {
			return nil, err
		}
		for k, v := range secrets {
			result[k] = v
		}
	}
	if len(result) == 0 {
		return nil, ErrEmptySecrets
	}
	return result, nil
}
