package cmd

import (
	"github.com/spf13/cobra"
	"github.com/your-org/vaultpipe/internal/retry"
)

// retryConfig holds parsed retry flags shared across commands.
var retryConfig retry.Config

// registerRetryFlags attaches retry-related flags to the given command.
func registerRetryFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(
		&retryConfig.MaxAttempts,
		"retry-attempts",
		3,
		"number of times to retry a failed Vault request",
	)
	cmd.Flags().DurationVar(
		&retryConfig.Delay,
		"retry-delay",
		retry.DefaultConfig().Delay,
		"initial delay between retry attempts (e.g. 500ms, 1s)",
	)
	cmd.Flags().Float64Var(
		&retryConfig.Multiplier,
		"retry-multiplier",
		retry.DefaultConfig().Multiplier,
		"backoff multiplier applied to delay on each retry (1.0 = constant)",
	)
}

// buildRetryDoer constructs a retry.Doer from the current retryConfig.
func buildRetryDoer() *retry.Doer {
	return retry.New(retryConfig)
}
