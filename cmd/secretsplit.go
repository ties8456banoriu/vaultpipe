package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vaultpipe/vaultpipe/internal/secretsplit"
)

var secretsplitCmd = &cobra.Command{
	Use:   "secretsplit",
	Short: "Partition secrets into named groups by key prefix",
	RunE:  runSecretsplit,
}

func init() {
	secretSplitCmd := secretsplitCmd
	secretSplitCmd.Flags().StringSliceP("rules", "r", nil, "Split rules in name:PREFIX_ format (e.g. db:DB_,aws:AWS_)")
	secretSplitCmd.Flags().StringP("input", "i", "", "Input .env file (default: stdin)")
	secretSplitCmd.Flags().StringP("output", "o", "json", "Output format: json or dotenv")
	rootCmd.AddCommand(secretSplitCmd)
}

func runSecretsplit(cmd *cobra.Command, _ []string) error {
	rulesRaw, _ := cmd.Flags().GetStringSlice("rules")
	inputFile, _ := cmd.Flags().GetString("input")
	outFmt, _ := cmd.Flags().GetString("output")

	var rules []secretsplit.Rule
	for _, r := range rulesRaw {
		parts := strings.SplitN(r, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid rule %q: expected name:PREFIX format", r)
		}
		rules = append(rules, secretsplit.Rule{Name: parts[0], Prefix: parts[1]})
	}

	splitter, err := secretsplit.New(rules)
	if err != nil {
		return fmt.Errorf("secretsplit: %w", err)
	}

	secrets, err := loadSecretsFromFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	groups, err := splitter.Apply(secrets)
	if err != nil {
		return fmt.Errorf("splitting secrets: %w", err)
	}

	switch outFmt {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(groups)
	case "dotenv":
		for group, kv := range groups {
			fmt.Fprintf(os.Stdout, "# group: %s\n", group)
			for k, v := range kv {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown output format %q", outFmt)
	}
}

func loadSecretsFromFile(path string) (map[string]string, error) {
	var data []byte
	var err error
	if path == "" {
		data, err = os.ReadFile("/dev/stdin")
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, err
	}
	out := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			out[parts[0]] = parts[1]
		}
	}
	return out, nil
}
