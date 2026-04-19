// Package secretencode provides base64 and base64url encoding and decoding
// transformations for secret values fetched from HashiCorp Vault.
//
// Use New to create an Encoder with the desired encoding type and direction,
// then call Apply to transform a map of secret key-value pairs.
package secretencode
