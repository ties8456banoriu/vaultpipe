package profile_test

import "errors"

// This file exists solely to make the errors import available to profile_test.go
// without embedding it in the main test file.
var _ = errors.New
