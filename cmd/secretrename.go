package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpipe/internal/secretrename"
)

var secretrenameRules []string

func init() {
	secretnameCmd := &cobra.Command{
		Use:   "secretrename",
		Short: "Rename secret keys using FROM:TO rules",
		Long: `Rename one or more secret keys in the current secrets map.

Rules are specified as FROM:TO pairs, e.g.:
  vaultpipe secretrename --rule DB_HOST:DATABASE_HOST --rule API_KEY:SERVICE_KEY`,
		RunE: runSecretrename,
	}
	secretnameCmd.Flags().StringArrayVar(&secretrenameRules, "rule", nil, "Rename rule in FROM:TO format (repeatable)")
	_ = secretnameCmd.MarkFlagRequired("rule")
	rootCmd.AddCommand(secretnameCmd)
}

func runSecretrename(cmd *cobra.Command, _ []string) error {
	rules, err := secretrename.ParseRules(secretrenameRules)
	if err != nil {
		return fmt.Errorf("parsing rename rules: %w", err)
	}

	renamer, err := secretrename.New(rules)
	if err != nil {
		return fmt.Errorf("creating renamer: %w", err)
	}

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	out, err := renamer.Apply(secrets)
	if err != nil {
		return fmt.Errorf("applying rename rules: %w", err)
	}

	for k, v := range out {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
	}
	return nil
}
