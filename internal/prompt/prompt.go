// Package prompt provides an interactive CLI prompt for selecting and
// confirming secret paths before injection into the local environment.
package prompt

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrNoSelection is returned when the user provides no input.
var ErrNoSelection = errors.New("prompt: no selection made")

// ErrCancelled is returned when the user explicitly cancels (enters 'q' or 'quit').
var ErrCancelled = errors.New("prompt: cancelled by user")

// Prompter presents choices to the user and returns the selected item.
type Prompter struct {
	in  io.Reader
	out io.Writer
}

// New creates a Prompter that reads from in and writes to out.
func New(in io.Reader, out io.Writer) *Prompter {
	return &Prompter{in: in, out: out}
}

// Choose displays a numbered list of options and returns the chosen item.
// Returns ErrCancelled if the user types 'q' or 'quit'.
// Returns ErrNoSelection if the input is empty.
func (p *Prompter) Choose(label string, options []string) (string, error) {
	if len(options) == 0 {
		return "", errors.New("prompt: no options provided")
	}

	fmt.Fprintf(p.out, "%s\n", label)
	for i, opt := range options {
		fmt.Fprintf(p.out, "  [%d] %s\n", i+1, opt)
	}
	fmt.Fprintf(p.out, "Enter number (or 'q' to cancel): ")

	var raw string
	buf := new(strings.Builder)
	tmp := make([]byte, 1)
	for {
		n, err := p.in.Read(tmp)
		if n > 0 {
			if tmp[0] == '\n' {
				break
			}
			buf.WriteByte(tmp[0])
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", fmt.Errorf("prompt: read error: %w", err)
		}
	}
	raw = strings.TrimSpace(buf.String())

	if raw == "" {
		return "", ErrNoSelection
	}
	if strings.EqualFold(raw, "q") || strings.EqualFold(raw, "quit") {
		return "", ErrCancelled
	}

	var idx int
	if _, err := fmt.Sscanf(raw, "%d", &idx); err != nil || idx < 1 || idx > len(options) {
		return "", fmt.Errorf("prompt: invalid selection %q", raw)
	}
	return options[idx-1], nil
}

// Confirm asks a yes/no question and returns true for 'y' or 'yes'.
func (p *Prompter) Confirm(question string) (bool, error) {
	fmt.Fprintf(p.out, "%s [y/N]: ", question)

	buf := new(strings.Builder)
	tmp := make([]byte, 1)
	for {
		n, err := p.in.Read(tmp)
		if n > 0 {
			if tmp[0] == '\n' {
				break
			}
			buf.WriteByte(tmp[0])
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return false, fmt.Errorf("prompt: read error: %w", err)
		}
	}
	answer := strings.TrimSpace(strings.ToLower(buf.String()))
	return answer == "y" || answer == "yes", nil
}
