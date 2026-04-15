package diff_test

import (
	"testing"

	"github.com/yourorg/vaultpipe/internal/diff"
)

func TestCompare_NoDifference(t *testing.T) {
	prev := map[string]string{"FOO": "bar", "BAZ": "qux"}
	curr := map[string]string{"FOO": "bar", "BAZ": "qux"}

	res := diff.Compare(prev, curr)
	if !res.IsEmpty() {
		t.Errorf("expected no diff, got %+v", res)
	}
}

func TestCompare_AddedKey(t *testing.T) {
	prev := map[string]string{"FOO": "bar"}
	curr := map[string]string{"FOO": "bar", "NEW": "val"}

	res := diff.Compare(prev, curr)
	if len(res.Added) != 1 || res.Added[0] != "NEW" {
		t.Errorf("expected Added=[NEW], got %v", res.Added)
	}
	if len(res.Removed) != 0 || len(res.Changed) != 0 {
		t.Errorf("unexpected removed/changed: %+v", res)
	}
}

func TestCompare_RemovedKey(t *testing.T) {
	prev := map[string]string{"FOO": "bar", "OLD": "gone"}
	curr := map[string]string{"FOO": "bar"}

	res := diff.Compare(prev, curr)
	if len(res.Removed) != 1 || res.Removed[0] != "OLD" {
		t.Errorf("expected Removed=[OLD], got %v", res.Removed)
	}
}

func TestCompare_ChangedKey(t *testing.T) {
	prev := map[string]string{"FOO": "old"}
	curr := map[string]string{"FOO": "new"}

	res := diff.Compare(prev, curr)
	if len(res.Changed) != 1 || res.Changed[0] != "FOO" {
		t.Errorf("expected Changed=[FOO], got %v", res.Changed)
	}
}

func TestCompare_MixedChanges(t *testing.T) {
	prev := map[string]string{"A": "1", "B": "2", "C": "3"}
	curr := map[string]string{"A": "1", "B": "changed", "D": "4"}

	res := diff.Compare(prev, curr)
	if len(res.Added) != 1 || res.Added[0] != "D" {
		t.Errorf("expected Added=[D], got %v", res.Added)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "C" {
		t.Errorf("expected Removed=[C], got %v", res.Removed)
	}
	if len(res.Changed) != 1 || res.Changed[0] != "B" {
		t.Errorf("expected Changed=[B], got %v", res.Changed)
	}
}

func TestCompare_EmptyBoth(t *testing.T) {
	res := diff.Compare(map[string]string{}, map[string]string{})
	if !res.IsEmpty() {
		t.Errorf("expected empty result for two empty maps, got %+v", res)
	}
}

func TestIsEmpty_False(t *testing.T) {
	res := diff.Result{Added: []string{"X"}}
	if res.IsEmpty() {
		t.Error("expected IsEmpty to return false")
	}
}
