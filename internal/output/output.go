// Package output formats and writes secret key-value pairs to various
// destinations such as stdout, a file, or a shell export block.
package output

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Format controls how secrets are rendered.
type Format string

const (
	FormatExport Format = "export" // export KEY=VALUE
	FormatDotenv Format = "dotenv" // KEY=VALUE
	FormatJSON   Format = "json"   // {"KEY":"VALUE"}
)

// ErrEmptySecrets is returned when no secrets are provided.
var ErrEmptySecrets = errors.New("output: secrets map is empty")

// Writer renders secrets to an io.Writer in the requested format.
type Writer struct {
	w      io.Writer
	format Format
}

// New creates a Writer that writes to w using the given format.
// If format is unrecognised it defaults to FormatDotenv.
func New(w io.Writer, format Format) *Writer {
	switch format {
	case FormatExport, FormatDotenv, FormatJSON:
	default:
		format = FormatDotenv
	}
	return &Writer{w: w, format: format}
}

// Write renders secrets to the underlying writer.
func (wr *Writer) Write(secrets map[string]string) error {
	if len(secrets) == 0 {
		return ErrEmptySecrets
	}
	switch wr.format {
	case FormatJSON:
		return wr.writeJSON(secrets)
	case FormatExport:
		return wr.writeLines(secrets, "export %s=%s\n")
	default:
		return wr.writeLines(secrets, "%s=%s\n")
	}
}

func (wr *Writer) writeLines(secrets map[string]string, tmpl string) error {
	for k, v := range secrets {
		if strings.ContainsAny(v, " \t") {
			v = `"` + v + `"`
		}
		if _, err := fmt.Fprintf(wr.w, tmpl, k, v); err != nil {
			return fmt.Errorf("output: write line: %w", err)
		}
	}
	return nil
}

func (wr *Writer) writeJSON(secrets map[string]string) error {
	pairs := make([]string, 0, len(secrets))
	for k, v := range secrets {
		pairs = append(pairs, fmt.Sprintf(`%q:%q`, k, v))
	}
	_, err := fmt.Fprintf(wr.w, "{%s}\n", strings.Join(pairs, ","))
	if err != nil {
		return fmt.Errorf("output: write json: %w", err)
	}
	return nil
}
