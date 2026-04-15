package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"vaultpipe/internal/audit"
	"vaultpipe/internal/config"
	"vaultpipe/internal/envwriter"
	"vaultpipe/internal/filter"
	"vaultpipe/internal/mapping"
	"vaultpipe/internal/vault"
)

var (
	filterRules []string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Fetch secrets from Vault and write them to a .env file",
	RunE:  runRun,
}

func init() {
	runCmd.Flags().StringArrayVar(&filterRules, "filter", nil, "Include/exclude rules for secret keys (e.g. DB_*, !DB_PASSWORD)")
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("config validation: %w", err)
	}

	auditLog, err := audit.NewLogger(cfg.AuditLogPath)
	if err != nil {
		return fmt.Errorf("audit logger: %w", err)
	}
	defer auditLog.Close()

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := client.GetSecret(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("vault fetch: %w", err)
	}
	auditLog.LogSecretFetched(cfg.SecretPath, len(secrets))

	if len(filterRules) > 0 {
		rules, err := filter.ParseRules(filterRules)
		if err != nil {
			return fmt.Errorf("filter rules: %w", err)
		}
		f := filter.NewFilter(rules)
		secrets, err = f.Apply(secrets)
		if err != nil {
			return fmt.Errorf("filter apply: %w", err)
		}
		log.Printf("filter: %d keys after filtering", len(secrets))
	}

	mapper := mapping.NewMapper(cfg.MappingRules)
	mapped, err := mapper.Apply(secrets)
	if err != nil {
		return fmt.Errorf("mapping: %w", err)
	}

	writer := envwriter.NewWriter(cfg.EnvFilePath)
	if err := writer.Write(mapped); err != nil {
		return fmt.Errorf("env write: %w", err)
	}
	auditLog.LogEnvWritten(cfg.EnvFilePath, len(mapped))

	fmt.Printf("vaultpipe: wrote %d secrets to %s\n", len(mapped), cfg.EnvFilePath)
	return nil
}
