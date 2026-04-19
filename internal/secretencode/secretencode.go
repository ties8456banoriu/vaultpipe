// Package secretencode provides encoding and decoding transformations for secret values.
package secretencode

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// Encoding represents a supported encoding type.
type Encoding string

const (
	EncodingBase64    Encoding = "base64"
	EncodingBase64URL Encoding = "base64url"
	EncodingHex       Encoding = "hex"
)

// Encoder applies encoding or decoding to secret values.
type Encoder struct {
	encoding Encoding
	decode   bool
}

// New creates an Encoder. Set decode=true to decode instead of encode.
func New(enc Encoding, decode bool) (*Encoder, error) {
	switch enc {
	case EncodingBase64, EncodingBase64URL:
		// valid
	default:
		return nil, fmt.Errorf("secretencode: unsupported encoding %q", enc)
	}
	return &Encoder{encoding: enc, decode: decode}, nil
}

// Apply encodes or decodes all values in secrets.
func (e *Encoder) Apply(secrets map[string]string) (map[string]string, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretencode: secrets map is empty")
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result, err := e.transform(v)
		if err != nil {
			return nil, fmt.Errorf("secretencode: key %q: %w", k, err)
		}
		out[k] = result
	}
	return out, nil
}

func (e *Encoder) transform(value string) (string, error) {
	if e.decode {
		return e.decodeValue(value)
	}
	return e.encodeValue(value), nil
}

func (e *Encoder) encodeValue(value string) string {
	switch e.encoding {
	case EncodingBase64URL:
		return base64.URLEncoding.EncodeToString([]byte(value))
	default:
		return base64.StdEncoding.EncodeToString([]byte(value))
	}
}

func (e *Encoder) decodeValue(value string) (string, error) {
	value = strings.TrimSpace(value)
	var b []byte
	var err error
	switch e.encoding {
	case EncodingBase64URL:
		b, err = base64.URLEncoding.DecodeString(value)
	default:
		b, err = base64.StdEncoding.DecodeString(value)
	}
	if err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}
	return string(b), nil
}
