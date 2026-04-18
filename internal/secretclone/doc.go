// Package secretclone provides a Cloner that produces deep copies of secret
// maps with optional key transformations such as prefix injection and
// uppercasing. It is useful when secrets need to be handed off to multiple
// pipeline stages without risk of cross-stage mutation.
package secretclone
