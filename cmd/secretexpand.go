package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secretexpand"
)

func init() {
	var secretexpandCmd = &cobra.Command{
		Use:   "secretexpand",
		Short: "Preview interpolated secret values from a .env file",
		Long: `Reads KEY=VALUE pairs from stdin or a file and resolves
${KEY} references within values, printing the expanded result.`,
		RunE: runSecretexpand,
	}
	rootCmd.AddCommand(secretexpandCmd)
}

func runSecretexpand(cmd *cobra.Command, args []string) error {
	// For demonstration: expand a hardcoded example map.
	// In production this would read from a .env file or vault fetch.
	example := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"DB_DSN":  "postgres://${DB_HOST}:${DB_PORT}/app",
	}

	e := secretexpand.New()
	result, err := e.Apply(example)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}

	for k, v := range result {
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}
