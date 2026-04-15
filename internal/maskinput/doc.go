// Package maskinput provides a [Prompter] for reading sensitive values such as
// Vault tokens from the terminal without echoing the characters back to the
// screen.
//
// When the underlying file descriptor is a real TTY, [Prompter.Prompt] uses
// golang.org/x/term.ReadPassword to suppress echo. When the descriptor is a
// pipe (e.g. in tests or scripted usage) input is read as plain text so that
// non-interactive workflows are still supported.
package maskinput
