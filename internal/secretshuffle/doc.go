// Package secretshuffle randomises the key ordering of a secrets map.
//
// This is useful when you need to process or display secrets in a
// non-deterministic order to avoid leaking information through consistent
// ordering patterns.
//
// Example:
//
//	s := secretshuffle.New(nil)
//	keys, err := s.Shuffle(secrets)
package secretshuffle
