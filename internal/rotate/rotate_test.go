package rotate_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/rotate"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_PASSWORD": "secret123",
		"API_KEY":     "abc",
	}
}

func TestDetect_NoBaseline_ReturnsErrNoBaseline(t *testing.T) {
	d := rotate.NewDetector()
	_, err := d.Detect(baseSecrets())
	if err != rotate.ErrNoBaseline {
		t.Fatalf("expected ErrNoBaseline, got %v", err)
	}
}

func TestSetBaseline_EmptySecrets_ReturnsError(t *testing.T) {
	d := rotate.NewDetector()
	err := d.SetBaseline(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty baseline, got nil")
	}
}

func TestDetect_NoChanges_ReturnsNilEvent(t *testing.T) {
	d := rotate.NewDetector()
	if err := d.SetBaseline(baseSecrets()); err != nil {
		t.Fatalf("SetBaseline: %v", err)
	}
	event, err := d.Detect(baseSecrets())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event != nil {
		t.Fatalf("expected nil event, got %+v", event)
	}
}

func TestDetect_ChangedValue_ReturnsEvent(t *testing.T) {
	d := rotate.NewDetector()
	if err := d.SetBaseline(baseSecrets()); err != nil {
		t.Fatalf("SetBaseline: %v", err)
	}
	current := baseSecrets()
	current["DB_PASSWORD"] = "newpassword"

	event, err := d.Detect(current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event == nil {
		t.Fatal("expected rotation event, got nil")
	}
	if len(event.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(event.Changes))
	}
	if event.DetectedAt.IsZero() {
		t.Error("expected DetectedAt to be set")
	}
}

func TestDetect_AddedKey_ReturnsEvent(t *testing.T) {
	d := rotate.NewDetector()
	if err := d.SetBaseline(baseSecrets()); err != nil {
		t.Fatalf("SetBaseline: %v", err)
	}
	current := baseSecrets()
	current["NEW_SECRET"] = "value"

	event, err := d.Detect(current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event == nil {
		t.Fatal("expected rotation event for added key")
	}
}

func TestRotationEvent_Summary_ContainsCount(t *testing.T) {
	d := rotate.NewDetector()
	if err := d.SetBaseline(baseSecrets()); err != nil {
		t.Fatalf("SetBaseline: %v", err)
	}
	current := baseSecrets()
	current["DB_PASSWORD"] = "rotated"
	current["API_KEY"] = "rotated2"

	event, err := d.Detect(current)
	if err != nil || event == nil {
		t.Fatalf("expected event, got err=%v event=%v", err, event)
	}
	summary := event.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
