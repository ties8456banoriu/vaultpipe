// Package masksecret provides middleware that redacts secret values from
// log output and written .env files based on a configured policy.
package masksecret

import (
	"errors"
	"strings"
)

// Policy controls which keys are masked in output.
type Policy int

const (
	// PolicyAll masks every value.
	PolicyAll Policy = iota
	// PolicyNone masks nothing (useful for debugging).
	PolicyNone
	// PolicySensitive masks only keys whose names contain a sensitive keyword.
	PolicySensitive
)

var sensitiveKeywords = []string{"password", "secret", "token", "key", "credential", "passwd", "apikey"}

// Masker applies masking to a map of secret key/value pairs.
type Masker struct {
	policy  Policy
	placeholder string
}

// New returns a Masker with the given policy. placeholder is the string
// substituted for masked values; if empty, "***" is used.
func New(policy Policy, placeholder string) (*Masker, error) {
	if placeholder == "" {
		placeholder = "***"
	}
	return &Masker{policy: policy, placeholder: placeholder}, nil
}

// Apply returns a copy of secrets with values masked according to the policy.
func (m *Masker) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("masksecret: secrets map is empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		switch m.policy {
		case PolicyAll:
			out[k] = m.placeholder
		case PolicyNone:
			out[k] = v
		case PolicySensitive:
			if isSensitive(k) {
				out[k] = m.placeholder
			} else {
				out[k] = v
			}
		}
	}
	return out, nil
}

func isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, kw := range sensitiveKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}
