// Package secretmergetag provides utilities for merging tag maps
// across multiple secret sets with configurable conflict resolution.
//
// Supported conflict policies:
//   - skip:      keep the first value on conflict
//   - overwrite: use the latest value on conflict
//   - error:     return an error on any conflict
package secretmergetag
