package secretscore_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretscore"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_PASS":  "Secr3tPass",
		"API_KEY":  "abc",
		"EMPTY_VAL": "",
	}
}

func TestApply_EmptySecrets_ReturnsError(t *testing.T) {
	s := secretscore.New()
	_, err := s.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestApply_ReturnsOneScorePerKey(t *testing.T) {
	s := secretscore.New()
	results, err := s.Apply(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 scores, got %d", len(results))
	}
}

func TestApply_PerfectScore(t *testing.T) {
	s := secretscore.New(secretscore.WithMinLength(8), secretscore.WithRequireUpper(true), secretscore.WithRequireDigit(true))
	results, err := s.Apply(map[string]string{"KEY": "Secr3tPass"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Points != results[0].MaxPoints {
		t.Errorf("expected perfect score, got %d/%d reasons: %v", results[0].Points, results[0].MaxPoints, results[0].Reasons)
	}
}

func TestApply_EmptyValue_ZeroPoints(t *testing.T) {
	s := secretscore.New()
	results, err := s.Apply(map[string]string{"KEY": ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Points != 0 {
		t.Errorf("expected 0 points for empty value, got %d", results[0].Points)
	}
	if len(results[0].Reasons) == 0 {
		t.Error("expected at least one reason for empty value")
	}
}

func TestApply_ShortValue_MissesLengthPoint(t *testing.T) {
	s := secretscore.New(secretscore.WithMinLength(16))
	results, err := s.Apply(map[string]string{"KEY": "Ab1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, r := range results[0].Reasons {
		if r == "value too short" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'value too short' reason")
	}
}

func TestPercent_FullScore_Returns100(t *testing.T) {
	sc := secretscore.Score{Points: 4, MaxPoints: 4}
	if sc.Percent() != 100 {
		t.Errorf("expected 100, got %v", sc.Percent())
	}
}

func TestPercent_ZeroMaxPoints_ReturnsZero(t *testing.T) {
	sc := secretscore.Score{Points: 0, MaxPoints: 0}
	if sc.Percent() != 0 {
		t.Errorf("expected 0, got %v", sc.Percent())
	}
}

func TestWithRequireUpper_False_SkipsUpperCheck(t *testing.T) {
	s := secretscore.New(secretscore.WithRequireUpper(false), secretscore.WithRequireDigit(false))
	results, err := s.Apply(map[string]string{"KEY": "alllower1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results[0].Reasons {
		if r == "no uppercase letter" {
			t.Error("unexpected 'no uppercase letter' reason when disabled")
		}
	}
}
