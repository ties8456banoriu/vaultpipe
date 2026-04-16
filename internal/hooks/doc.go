// Package hooks provides lifecycle hook support for vaultpipe.
//
// Hooks allow users to run arbitrary shell commands before (pre) or after
// (post) secrets are fetched and written. Rules are specified as
// "stage:command" strings, e.g. "pre:make vault-login".
package hooks
