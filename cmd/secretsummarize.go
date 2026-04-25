package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wunderkind/vaultpipe/internal/secretsummarize"
)

func init() {
	var secretsummarizeCmd = &cobra.Command{
		Use:   "secret-summarize",
		Short: "Print a statistical summary of secrets loaded from a .env file",
		RunE:  runSecretsummarize,
	}

	secretSummarizeCmd := secretsummarizeCmd
	secretSummarizeCmd.Flags().StringP("file", "f", ".env", "Path to the .env file to summarize")
	rootCmd.AddCommand(secretSummarizeCmd)
}

func runSecretsummarize(cmd *cobra.Command, _ []string) error {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("secret-summarize: cannot read file %q: %w", filePath, err)
	}

	secrets := parseEnvLines(string(data))

	sum, err := secretsummarize.New().Apply(secrets)
	if err != nil {
		return fmt.Errorf("secret-summarize: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Total keys:     %d\n", sum.TotalKeys)
	fmt.Fprintf(cmd.OutOrStdout(), "Empty values:   %d\n", sum.EmptyValues)
	fmt.Fprintf(cmd.OutOrStdout(), "Unique values:  %d\n", sum.UniqueValues)
	fmt.Fprintf(cmd.OutOrStdout(), "Avg value len:  %.2f\n", sum.AvgValueLen)
	fmt.Fprintf(cmd.OutOrStdout(), "Longest key:    %s\n", sum.LongestKey)
	fmt.Fprintf(cmd.OutOrStdout(), "Shortest key:   %s\n", sum.ShortestKey)
	return nil
}

// parseEnvLines parses a simple KEY=VALUE .env file into a map.
func parseEnvLines(content string) map[string]string {
	result := make(map[string]string)
	for _, line := range splitLines(content) {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		for i, ch := range line {
			if ch == '=' {
				key := line[:i]
				val := line[i+1:]
				if key != "" {
					result[key] = val
				}
				break
			}
		}
	}
	return result
}
