// Package secretrotate tracks secret rotation events and policies for
// secrets fetched from HashiCorp Vault. It records when each secret was
// last rotated, which version was written, and whether rotation was
// triggered manually or on a schedule. Use NeedsRotation to determine
// whether a secret has exceeded its maximum allowed age.
package secretrotate
