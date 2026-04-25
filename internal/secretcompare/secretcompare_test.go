package secretcompare_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/secretcompare"
)

func TestCompare_BothEmpty_ReturnsError(t *testing.T) {
	c := secretcompare.New()
	_, err := c.Compare(nil, nil)
	if err == nil {
		t.Fatal("expected error for empty maps, got nil")
	}
}

func TestCompare_AllMatch(t *testing.T) {
	c := secretcompare.New()
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	results, err := c.Compare(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Status != secretcompare.StatusMatch {
			t.Errorf("key %q: expected match, got %s", r.Key, r.Status)
		}
	}
}

func TestCompare_Mismatch(t *testing.T) {
	c := secretcompare.New()
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "baz"}
	results, err := c.Compare(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != secretcompare.StatusMismatch {
		t.Errorf("expected mismatch, got %s", results[0].Status)
	}
}

func TestCompare_MissingInB(t *testing.T) {
	c := secretcompare.New()
	a := map[string]string{"FOO": "bar", "ONLY_A": "x"}
	b := map[string]string{"FOO": "bar"}
	results, err := c.Compare(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "ONLY_A" && r.Status != secretcompare.StatusMissingB {
			t.Errorf("expected missing_b for ONLY_A, got %s", r.Status)
		}
	}
}

func TestCompare_MissingInA(t *testing.T) {
	c := secretcompare.New()
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "y"}
	results, err := c.Compare(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "ONLY_B" && r.Status != secretcompare.StatusMissingA {
			t.Errorf("expected missing_a for ONLY_B, got %s", r.Status)
		}
	}
}

func TestCompare_ResultsAreSorted(t *testing.T) {
	c := secretcompare.New()
	a := map[string]string{"ZZZ": "1", "AAA": "2", "MMM": "3"}
	b := map[string]string{"ZZZ": "1", "AAA": "2", "MMM": "3"}
	results, err := c.Compare(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"AAA", "MMM", "ZZZ"}
	for i, r := range results {
		if r.Key != expected[i] {
			t.Errorf("position %d: expected key %q, got %q", i, expected[i], r.Key)
		}
	}
}

func TestSummary_CorrectCounts(t *testing.T) {
	results := []secretcompare.Result{
		{Key: "A", Status: secretcompare.StatusMatch},
		{Key: "B", Status: secretcompare.StatusMismatch},
		{Key: "C", Status: secretcompare.StatusMissingA},
		{Key: "D", Status: secretcompare.StatusMissingB},
	}
	got := secretcompare.Summary(results)
	want := "match=1 mismatch=1 missing_a=1 missing_b=1"
	if got != want {
		t.Errorf("Summary: got %q, want %q", got, want)
	}
}
