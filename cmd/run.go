package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/vaultpipe/internal/config"
	"github.com/yourorg/vaultpipe/internal/envwriter"
	"github.com/yourorg/vaultpipe/internal/mapping"
	"github.com/yourorg/vaultpipe/internal/vault"
)

var mapRules []string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Fetch secrets from Vault and write them to a .env file",
	RunE:  runRun,
}

func init() {
	runCmd.Flags().StringArrayVar(&mapRules, "map", []string{}, "Key mapping rules in vault_key=ENV_KEY format")
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	secrets, err := client.GetSecret(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("fetching secret %q: %w", cfg.SecretPath, err)
	}

	rules, err := mapping.ParseRules(mapRules)
	if err != nil {
		return fmt.Errorf("parsing map rules: %w", err)
	}

	mapper := mapping.NewMapper(rules)
	mapped, err := mapper.Apply(secrets)
	if err != nil {
		return fmt.Errorf("applying mappings: %w", err)
	}

	w, err := envwriter.NewWriter(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("creating env writer: %w", err)
	}

	if err := w.Write(mapped); err != nil {
		return fmt.Errorf("writing .env file: %w", err)
	}

	fmt.Fprintf(os.Stdout, "vaultpipe: wrote %d secrets to %s\n", len(mapped), cfg.OutputFile)
	return nil
}
