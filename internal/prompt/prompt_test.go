package prompt

import (
	"bytes"
	"strings"
	"testing"
)

func makePrompter(input string) (*Prompter, *bytes.Buffer) {
	out := &bytes.Buffer{}
	return New(strings.NewReader(input), out), out
}

func TestChoose_ValidSelection(t *testing.T) {
	p, _ := makePrompter("2\n")
	options := []string{"secret/dev", "secret/staging", "secret/prod"}
	got, err := p.Choose("Select a path:", options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/staging" {
		t.Errorf("expected secret/staging, got %q", got)
	}
}

func TestChoose_CancelWithQ(t *testing.T) {
	p, _ := makePrompter("q\n")
	_, err := p.Choose("Select:", []string{"a", "b"})
	if err != ErrCancelled {
		t.Errorf("expected ErrCancelled, got %v", err)
	}
}

func TestChoose_EmptyInput_ReturnsErrNoSelection(t *testing.T) {
	p, _ := makePrompter("\n")
	_, err := p.Choose("Select:", []string{"a", "b"})
	if err != ErrNoSelection {
		t.Errorf("expected ErrNoSelection, got %v", err)
	}
}

func TestChoose_OutOfRange_ReturnsError(t *testing.T) {
	p, _ := makePrompter("99\n")
	_, err := p.Choose("Select:", []string{"a", "b"})
	if err == nil {
		t.Fatal("expected error for out-of-range selection")
	}
}

func TestChoose_NoOptions_ReturnsError(t *testing.T) {
	p, _ := makePrompter("1\n")
	_, err := p.Choose("Select:", []string{})
	if err == nil {
		t.Fatal("expected error for empty options")
	}
}

func TestChoose_PrintsLabel(t *testing.T) {
	p, out := makePrompter("1\n")
	_, _ = p.Choose("Pick secret path:", []string{"secret/dev"})
	if !strings.Contains(out.String(), "Pick secret path:") {
		t.Errorf("expected label in output, got: %q", out.String())
	}
}

func TestConfirm_Yes(t *testing.T) {
	p, _ := makePrompter("y\n")
	ok, err := p.Confirm("Overwrite .env?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true for 'y'")
	}
}

func TestConfirm_No(t *testing.T) {
	p, _ := makePrompter("n\n")
	ok, err := p.Confirm("Overwrite .env?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false for 'n'")
	}
}

func TestConfirm_EmptyInput_ReturnsFalse(t *testing.T) {
	p, _ := makePrompter("\n")
	ok, err := p.Confirm("Continue?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false for empty input")
	}
}
