package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secretcopy"
)

var secretcopyCmdRules []string

func init() {
	var cmd = &cobra.Command{
		Use:   "secretcopy",
		Short: "Copy secrets from one key to another using FROM:TO rules",
		Long: `Copy secret values to new keys within the loaded secret map.
Rules are specified as FROM:TO pairs. Original keys are preserved.`,
		RunE: runSecretcopy,
	}
	cmd.Flags().StringArrayVar(&secretcopyCmdRules, "rule", nil, "Copy rule in FROM:TO format (repeatable)")
	_ = cmd.MarkFlagRequired("rule")
	rootCmd.AddCommand(cmd)
}

func runSecretcopy(cmd *cobra.Command, _ []string) error {
	rules, err := secretcopy.ParseRules(secretcopyCmdRules)
	if err != nil {
		return fmt.Errorf("invalid rules: %w", err)
	}

	copier, err := secretcopy.New(rules)
	if err != nil {
		return fmt.Errorf("failed to create copier: %w", err)
	}

	// Placeholder: in production this would be wired to the vault client.
	secrets := map[string]string{}
	out, err := copier.Apply(secrets)
	if err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	for k, v := range out {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
	}
	return nil
}
