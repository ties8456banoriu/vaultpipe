// Package redact provides utilities for redacting sensitive secret values
// before they are displayed in logs, terminal output, or error messages.
package redact

import "strings"

const (
	// DefaultMaskChar is the character used to mask secret values.
	DefaultMaskChar = "*"
	// DefaultVisibleChars is the number of trailing characters to reveal.
	DefaultVisibleChars = 4
)

// Redactor masks secret values to prevent accidental exposure.
type Redactor struct {
	maskChar     string
	visibleChars int
}

// NewRedactor creates a Redactor with default settings.
func NewRedactor() *Redactor {
	return &Redactor{
		maskChar:     DefaultMaskChar,
		visibleChars: DefaultVisibleChars,
	}
}

// Mask replaces most of the value with mask characters, revealing only the
// last visibleChars characters. Short values are fully masked.
func (r *Redactor) Mask(value string) string {
	if len(value) == 0 {
		return ""
	}
	if len(value) <= r.visibleChars {
		return strings.Repeat(r.maskChar, len(value))
	}
	maskLen := len(value) - r.visibleChars
	return strings.Repeat(r.maskChar, maskLen) + value[maskLen:]
}

// MaskAll replaces every character in value with the mask character.
func (r *Redactor) MaskAll(value string) string {
	if len(value) == 0 {
		return ""
	}
	return strings.Repeat(r.maskChar, len(value))
}

// MaskMap returns a new map where all values are masked.
// Keys are left unchanged.
func (r *Redactor) MaskMap(secrets map[string]string) map[string]string {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = r.Mask(v)
	}
	return result
}
