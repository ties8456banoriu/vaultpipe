// Package labels provides tagging support for secret env keys.
//
// Labels are arbitrary key-value pairs (e.g. env=production, team=backend)
// that can be attached to individual secrets. They can be used to filter
// which secrets are written, logged, or processed by downstream pipeline
// stages.
package labels
