package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secretpatch"
)

func init() {
	var policy string
	var keys []string

	cmd := &cobra.Command{
		Use:   "secretpatch",
		Short: "Patch base secrets with values from a patch set",
		Long: `Merges selected keys from a patch secrets map into a base secrets map.

Conflict resolution is controlled by --policy:
  overwrite  patch value replaces base value (default)
  keep       base value is preserved on conflict`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecretpatch(policy, keys)
		},
	}

	cmd.Flags().StringVar(&policy, "policy", "overwrite", "conflict resolution policy: overwrite|keep")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "comma-separated list of patch keys to apply (default: all)")

	rootCmd.AddCommand(cmd)
}

func runSecretpatch(policy string, keys []string) error {
	// Example wiring: in a real flow these would come from vault + a secondary path.
	base := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	patch := map[string]string{
		"DB_PORT": "5433",
		"DB_NAME": "mydb",
	}

	p, err := secretpatch.New(secretpatch.Policy(policy))
	if err != nil {
		return fmt.Errorf("secretpatch: %w", err)
	}

	result, err := p.Apply(base, patch, keys)
	if err != nil {
		return fmt.Errorf("secretpatch: %w", err)
	}

	for k, v := range result {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
	}
	return nil
}
