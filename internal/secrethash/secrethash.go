// Package secrethash provides hashing of secret values using configurable
// algorithms. It is useful for producing stable, non-reversible fingerprints
// of secret values for comparison or audit purposes.
package secrethash

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
)

// Algorithm represents a supported hashing algorithm.
type Algorithm string

const (
	AlgoSHA256 Algorithm = "sha256"
	AlgoSHA512 Algorithm = "sha512"
	AlgoMD5    Algorithm = "md5"
)

// Hasher hashes secret values in place.
type Hasher struct {
	algo Algorithm
}

// New returns a Hasher for the given algorithm.
// Returns an error if the algorithm is not supported.
func New(algo Algorithm) (*Hasher, error) {
	switch algo {
	case AlgoSHA256, AlgoSHA512, AlgoMD5:
		return &Hasher{algo: algo}, nil
	default:
		return nil, fmt.Errorf("secrethash: unsupported algorithm %q", algo)
	}
}

// Apply replaces each secret value with its hex-encoded hash.
// Returns an error if secrets is empty.
func (h *Hasher) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secrethash: secrets must not be empty")
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		hashed, err := h.hash(v)
		if err != nil {
			return nil, fmt.Errorf("secrethash: failed to hash key %q: %w", k, err)
		}
		out[k] = hashed
	}
	return out, nil
}

// Hash returns the hex-encoded hash of a single value.
func (h *Hasher) Hash(value string) (string, error) {
	return h.hash(value)
}

func (h *Hasher) hash(value string) (string, error) {
	var hw hash.Hash
	switch h.algo {
	case AlgoSHA256:
		hw = sha256.New()
	case AlgoSHA512:
		hw = sha512.New()
	case AlgoMD5:
		hw = md5.New()
	default:
		return "", fmt.Errorf("unsupported algorithm %q", h.algo)
	}
	_, err := hw.Write([]byte(value))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hw.Sum(nil)), nil
}
