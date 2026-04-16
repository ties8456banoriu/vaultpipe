package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/hooks"
)

// registerHookFlags adds --hook flags to a command.
func registerHookFlags(cmd *cobra.Command) {
	cmd.Flags().StringArray("hook", nil, `lifecycle hook in "stage:command" format (e.g. "pre:make login").
May be specified multiple times. Stages: pre, post.`)
}

// buildHookRunner parses hook flags from the command and returns a Runner.
func buildHookRunner(cmd *cobra.Command) (*hooks.Runner, error) {
	raw, err := cmd.Flags().GetStringArray("hook")
	if err != nil {
		return nil, err
	}
	parsed, err := hooks.ParseRules(raw)
	if err != nil {
		return nil, err
	}
	return hooks.NewRunner(parsed), nil
}
