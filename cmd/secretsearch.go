package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secretsearch"
)

func init() {
	var mode string
	var searchKey, searchVal bool

	cmd := &cobra.Command{
		Use:   "secret-search <query>",
		Short: "Search secrets by key or value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecretSearch(args[0], mode, searchKey, searchVal)
		},
	}

	cmd.Flags().StringVar(&mode, "mode", "prefix", "Search mode: exact, prefix, or regex")
	cmd.Flags().BoolVar(&searchKey, "key", true, "Match against secret keys")
	cmd.Flags().BoolVar(&searchVal, "value", false, "Match against secret values")

	rootCmd.AddCommand(cmd)
}

func runSecretSearch(query, mode string, searchKey, searchVal bool) error {
	secrets, err := loadSecretsFromEnvFile()
	if err != nil {
		return fmt.Errorf("secret-search: failed to load secrets: %w", err)
	}

	s, err := secretsearch.New(secretsearch.Mode(mode), searchKey, searchVal)
	if err != nil {
		return err
	}

	results, err := s.Search(secrets, query)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "no matches found")
		return nil
	}

	for _, r := range results {
		fmt.Fprintf(os.Stdout, "%s=%s\n", r.Key, r.Value)
	}
	return nil
}

// loadSecretsFromEnvFile reads KEY=VALUE lines from the configured .env file.
func loadSecretsFromEnvFile() (map[string]string, error) {
	path := os.Getenv("VAULTPIPE_ENV_FILE")
	if path == "" {
		path = ".env"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	secrets := make(map[string]string)
	for _, line := range splitLines(string(data)) {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		for i, ch := range line {
			if ch == '=' {
				secrets[line[:i]] = line[i+1:]
				break
			}
		}
	}
	return secrets, nil
}
