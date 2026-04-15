// Package profile manages named environment profiles, allowing users to
// define multiple Vault secret configurations (e.g. "dev", "staging") and
// switch between them easily.
package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrNoProfiles is returned when no profiles have been defined.
var ErrNoProfiles = errors.New("no profiles defined")

// ErrProfileNotFound is returned when a named profile does not exist.
var ErrProfileNotFound = errors.New("profile not found")

// Profile holds configuration for a named environment profile.
type Profile struct {
	Name       string            `json:"name"`
	SecretPath string            `json:"secret_path"`
	EnvFile    string            `json:"env_file"`
	Meta       map[string]string `json:"meta,omitempty"`
}

// Store holds a collection of named profiles persisted to disk.
type Store struct {
	path     string
	profiles map[string]Profile
}

// NewStore creates a Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path, profiles: make(map[string]Profile)}
}

// Load reads profiles from disk. Returns ErrNoProfiles if the file is absent.
func (s *Store) Load() error {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return ErrNoProfiles
	}
	if err != nil {
		return fmt.Errorf("profile: read file: %w", err)
	}
	var profiles []Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return fmt.Errorf("profile: parse file: %w", err)
	}
	for _, p := range profiles {
		s.profiles[p.Name] = p
	}
	return nil
}

// Save persists all profiles to disk, creating parent directories as needed.
func (s *Store) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o750); err != nil {
		return fmt.Errorf("profile: mkdir: %w", err)
	}
	list := make([]Profile, 0, len(s.profiles))
	for _, p := range s.profiles {
		list = append(list, p)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("profile: marshal: %w", err)
	}
	return os.WriteFile(s.path, data, 0o640)
}

// Set adds or replaces a profile.
func (s *Store) Set(p Profile) error {
	if p.Name == "" {
		return errors.New("profile: name must not be empty")
	}
	if p.SecretPath == "" {
		return errors.New("profile: secret_path must not be empty")
	}
	s.profiles[p.Name] = p
	return nil
}

// Get retrieves a profile by name.
func (s *Store) Get(name string) (Profile, error) {
	p, ok := s.profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("%w: %s", ErrProfileNotFound, name)
	}
	return p, nil
}

// List returns all stored profiles.
func (s *Store) List() []Profile {
	list := make([]Profile, 0, len(s.profiles))
	for _, p := range s.profiles {
		list = append(list, p)
	}
	return list
}

// Delete removes a profile by name. Returns ErrProfileNotFound if absent.
func (s *Store) Delete(name string) error {
	if _, ok := s.profiles[name]; !ok {
		return fmt.Errorf("%w: %s", ErrProfileNotFound, name)
	}
	delete(s.profiles, name)
	return nil
}
