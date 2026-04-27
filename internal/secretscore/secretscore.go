// Package secretscore computes a quality score for secrets based on
// configurable criteria such as value length, entropy, and key naming.
package secretscore

import (
	"errors"
	"math"
	"strings"
)

// Score holds the result of scoring a single secret.
type Score struct {
	Key      string
	Value    string
	Points   int
	MaxPoints int
	Reasons  []string
}

// Percent returns the score as a percentage (0–100).
func (s Score) Percent() float64 {
	if s.MaxPoints == 0 {
		return 0
	}
	return math.Round(float64(s.Points) / float64(s.MaxPoints) * 100)
}

// Scorer scores a map of secrets.
type Scorer struct {
	minLength int
	requireUpper bool
	requireDigit bool
}

// Option configures a Scorer.
type Option func(*Scorer)

// WithMinLength sets the minimum value length for a full length score.
func WithMinLength(n int) Option {
	return func(s *Scorer) { s.minLength = n }
}

// WithRequireUpper requires at least one uppercase letter for a full score.
func WithRequireUpper(v bool) Option {
	return func(s *Scorer) { s.requireUpper = v }
}

// WithRequireDigit requires at least one digit for a full score.
func WithRequireDigit(v bool) Option {
	return func(s *Scorer) { s.requireDigit = v }
}

// New creates a Scorer with the supplied options.
func New(opts ...Option) *Scorer {
	s := &Scorer{minLength: 8, requireUpper: true, requireDigit: true}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply scores every secret in the map and returns a slice of Score results.
func (s *Scorer) Apply(secrets map[string]string) ([]Score, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretscore: no secrets to score")
	}
	results := make([]Score, 0, len(secrets))
	for k, v := range secrets {
		results = append(results, s.score(k, v))
	}
	return results, nil
}

func (s *Scorer) score(key, value string) Score {
	max := 1 // non-empty
	if s.minLength > 0 {
		max++
	}
	if s.requireUpper {
		max++
	}
	if s.requireDigit {
		max++
	}

	sc := Score{Key: key, Value: value, MaxPoints: max}

	if value == "" {
		sc.Reasons = append(sc.Reasons, "empty value")
		return sc
	}
	sc.Points++

	if s.minLength > 0 {
		if len(value) >= s.minLength {
			sc.Points++
		} else {
			sc.Reasons = append(sc.Reasons, "value too short")
		}
	}
	if s.requireUpper {
		if strings.ContainsAny(value, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			sc.Points++
		} else {
			sc.Reasons = append(sc.Reasons, "no uppercase letter")
		}
	}
	if s.requireDigit {
		if strings.ContainsAny(value, "0123456789") {
			sc.Points++
		} else {
			sc.Reasons = append(sc.Reasons, "no digit")
		}
	}
	return sc
}
