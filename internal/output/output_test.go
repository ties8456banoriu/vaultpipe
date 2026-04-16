package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/output"
)

func TestWrite_EmptySecrets_ReturnsError(t *testing.T) {
	w := output.New(&bytes.Buffer{}, output.FormatDotenv)
	if err := w.Write(map[string]string{}); err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestWrite_DotenvFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatDotenv)
	if err := w.Write(map[string]string{"FOO": "bar"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", buf.String())
	}
}

func TestWrite_ExportFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatExport)
	if err := w.Write(map[string]string{"BAR": "baz"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "export BAR=baz") {
		t.Errorf("expected export BAR=baz, got: %s", buf.String())
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatJSON)
	if err := w.Write(map[string]string{"K": "v"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.HasPrefix(got, "{") || !strings.Contains(got, "}") {
		t.Errorf("expected JSON object, got: %s", got)
	}
}

func TestWrite_QuotesValuesWithSpaces(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatDotenv)
	if err := w.Write(map[string]string{"MSG": "hello world"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"hello world"`) {
		t.Errorf("expected quoted value, got: %s", buf.String())
	}
}

func TestNew_UnknownFormat_DefaultsToDotenv(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.Format("xml"))
	if err := w.Write(map[string]string{"X": "1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "X=1") {
		t.Errorf("expected dotenv fallback, got: %s", buf.String())
	}
}
