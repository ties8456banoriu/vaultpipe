// Package secretinspect provides a read-only inspector for examining
// individual secret entries, returning structured metadata about each key.
package secretinspect

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrEmptySecrets is returned when an empty secrets map is provided.
var ErrEmptySecrets = errors.New("secretinspect: secrets map is empty")

// ErrUnknownKey is returned when the requested key does not exist.
var ErrUnknownKey = errors.New("secretinspect: key not found")

// Info holds inspection details for a single secret entry.
type Info struct {
	Key       string
	Value     string
	Length    int
	IsEmpty   bool
	HasSpaces bool
	Uppercase bool
}

// Inspector inspects secret entries.
type Inspector struct{}

// New returns a new Inspector.
func New() *Inspector {
	return &Inspector{}
}

// Inspect returns an Info struct for the given key in secrets.
func (i *Inspector) Inspect(secrets map[string]string, key string) (Info, error) {
	if len(secrets) == 0 {
		return Info{}, ErrEmptySecrets
	}
	v, ok := secrets[key]
	if !ok {
		return Info{}, fmt.Errorf("%w: %q", ErrUnknownKey, key)
	}
	return Info{
		Key:       key,
		Value:     v,
		Length:    len(v),
		IsEmpty:   v == "",
		HasSpaces: strings.Contains(v, " "),
		Uppercase: v == strings.ToUpper(v) && v != "",
	}, nil
}

// All returns a slice of Info for every key in secrets, sorted alphabetically.
func (i *Inspector) All(secrets map[string]string) ([]Info, error) {
	if len(secrets) == 0 {
		return nil, ErrEmptySecrets
	}
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make([]Info, 0, len(keys))
	for _, k := range keys {
		info, err := i.Inspect(secrets, k)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}
	return result, nil
}
