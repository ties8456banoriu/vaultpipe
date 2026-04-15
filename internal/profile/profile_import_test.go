package profile_test

import (
	"encoding/json"
	"errors"
)

// This file provides shared test helpers for the profile_test package.

// Sentinel errors used across profile import tests.
var (
	_ = errors.New

	errInvalidProfile = errors.New("invalid profile data")
	errMissingField   = errors.New("missing required field")
)

// mustMarshalJSON marshals v to JSON and panics if marshalling fails.
// Useful in test setup where errors are not expected.
func mustMarshalJSON(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic("mustMarshalJSON: " + err.Error())
	}
	return data
}
