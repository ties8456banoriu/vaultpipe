package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultpipe/vaultpipe/internal/config"
)

var (
	cfg        *config.Config
	secretPath string
	outputFile string
	overwrite  bool
)

var rootCmd = &cobra.Command{
	Use:   "vaultpipe",
	Short: "Inject Vault secrets into local .env files",
	Long: `vaultpipe reads secrets from HashiCorp Vault and writes them
into a local .env file for use in development environments.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg = config.LoadFromEnv()

		// CLI flags override env-based config
		if secretPath != "" {
			cfg.SecretPath = secretPath
		}
		if outputFile != "" {
			cfg.OutputFile = outputFile
		}
		cfg.Overwrite = overwrite

		return cfg.Validate()
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&secretPath, "path", "p", "",
		"Vault secret path (e.g. secret/data/myapp)",
	)
	rootCmd.PersistentFlags().StringVarP(
		&outputFile, "output", "o", "",
		"Output .env file path (default: .env)",
	)
	rootCmd.PersistentFlags().BoolVar(
		&overwrite, "overwrite", false,
		"Overwrite existing keys in the .env file",
	)
}
