// Package maskinput provides utilities for securely prompting users for
// sensitive input (e.g. Vault tokens) without echoing characters to the terminal.
package maskinput

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// ErrEmptyInput is returned when the user provides no input.
var ErrEmptyInput = errors.New("maskinput: input must not be empty")

// Prompter reads masked input from a terminal or falls back to a plain reader.
type Prompter struct {
	In  *os.File
	out io.Writer
}

// NewPrompter creates a Prompter that reads from in and writes the prompt to out.
func NewPrompter(in *os.File, out io.Writer) *Prompter {
	return &Prompter{In: in, out: out}
}

// Prompt displays msg to the user, reads a masked line of input, and returns
// the trimmed result. If the file descriptor is not a terminal the input is
// read as plain text so that pipes and tests still work.
func (p *Prompter) Prompt(msg string) (string, error) {
	fmt.Fprint(p.out, msg)

	var raw []byte
	var err error

	fd := int(p.In.Fd())
	if term.IsTerminal(fd) {
		raw, err = term.ReadPassword(fd)
		fmt.Fprintln(p.out) // move to next line after hidden input
	} else {
		buf := new(strings.Builder)
		_, err = fmt.Fscan(p.In, buf)
		raw = []byte(buf.String())
	}

	if err != nil {
		return "", fmt.Errorf("maskinput: reading input: %w", err)
	}

	value := strings.TrimSpace(string(raw))
	if value == "" {
		return "", ErrEmptyInput
	}

	return value, nil
}
