package maskinput_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/maskinput"
)

// writePipe returns an *os.File backed by a pipe whose read-end contains data.
func writePipe(t *testing.T, data string) *os.File {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	_, _ = w.WriteString(data)
	w.Close()
	t.Cleanup(func() { r.Close() })
	return r
}

func TestPrompt_ReturnsInput(t *testing.T) {
	r := writePipe(t, "s3cr3t\n")
	out := &bytes.Buffer{}
	p := maskinput.NewPrompter(r, out)

	got, err := p.Prompt("Token: ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "s3cr3t" {
		t.Errorf("got %q, want %q", got, "s3cr3t")
	}
	if !strings.Contains(out.String(), "Token: ") {
		t.Errorf("prompt message not written to output")
	}
}

func TestPrompt_TrimsWhitespace(t *testing.T) {
	r := writePipe(t, "  mytoken  \n")
	p := maskinput.NewPrompter(r, &bytes.Buffer{})

	got, err := p.Prompt("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "mytoken" {
		t.Errorf("got %q, want trimmed %q", got, "mytoken")
	}
}

func TestPrompt_EmptyInput_ReturnsError(t *testing.T) {
	r := writePipe(t, "   \n")
	p := maskinput.NewPrompter(r, &bytes.Buffer{})

	_, err := p.Prompt("")
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
	if err != maskinput.ErrEmptyInput {
		t.Errorf("got %v, want ErrEmptyInput", err)
	}
}

func TestPrompt_ClosedPipe_ReturnsError(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	w.Close() // close immediately — read will EOF with no data
	defer r.Close()

	p := maskinput.NewPrompter(r, &bytes.Buffer{})
	_, err = p.Prompt("")
	if err == nil {
		t.Fatal("expected error on closed pipe, got nil")
	}
}
