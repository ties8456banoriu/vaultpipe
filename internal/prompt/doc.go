// Package prompt provides interactive CLI prompts used by vaultpipe commands
// to let users select secret paths and confirm destructive operations before
// they are applied to the local environment.
//
// The prompts in this package are built on top of terminal-based UI primitives
// and are designed to degrade gracefully when running in non-interactive
// environments (e.g. CI pipelines or piped input). In such cases, callers
// should check for [ErrNotInteractive] and fall back to flag-based input.
package prompt
