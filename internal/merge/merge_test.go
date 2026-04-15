package merge_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/merge"
)

func TestMerge_NoSources_ReturnsError(t *testing.T) {
	m := merge.NewMerger(merge.StrategyFirst)
	_, err := m.Merge()
	if err != merge.ErrNoSources {
		t.Fatalf("expected ErrNoSources, got %v", err)
	}
}

func TestMerge_AllEmpty_ReturnsError(t *testing.T) {
	m := merge.NewMerger(merge.StrategyFirst)
	_, err := m.Merge(map[string]string{}, map[string]string{})
	if err != merge.ErrEmptySources {
		t.Fatalf("expected ErrEmptySources, got %v", err)
	}
}

func TestMerge_StrategyFirst_NoConflict(t *testing.T) {
	m := merge.NewMerger(merge.StrategyFirst)
	a := map[string]string{"FOO": "foo"}
	b := map[string]string{"BAR": "bar"}
	got, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "foo" || got["BAR"] != "bar" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestMerge_StrategyFirst_KeepsFirstOnConflict(t *testing.T) {
	m := merge.NewMerger(merge.StrategyFirst)
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}
	got, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "from-a" {
		t.Errorf("expected 'from-a', got %q", got["KEY"])
	}
}

func TestMerge_StrategyLast_KeepsLastOnConflict(t *testing.T) {
	m := merge.NewMerger(merge.StrategyLast)
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}
	got, err := m.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "from-b" {
		t.Errorf("expected 'from-b', got %q", got["KEY"])
	}
}

func TestMerge_SingleSource_ReturnsClone(t *testing.T) {
	m := merge.NewMerger(merge.StrategyFirst)
	src := map[string]string{"A": "1", "B": "2"}
	got, err := m.Merge(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got["A"] != "1" || got["B"] != "2" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestMerge_StrategyLast_MultipleSources(t *testing.T) {
	m := merge.NewMerger(merge.StrategyLast)
	a := map[string]string{"X": "a", "Y": "a"}
	b := map[string]string{"X": "b"}
	c := map[string]string{"X": "c", "Z": "c"}
	got, err := m.Merge(a, b, c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["X"] != "c" || got["Y"] != "a" || got["Z"] != "c" {
		t.Errorf("unexpected result: %v", got)
	}
}
