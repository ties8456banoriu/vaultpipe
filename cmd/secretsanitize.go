package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultpipe/vaultpipe/internal/secretsanitize"
)

func init() {
	var (
		trimSpace       bool
		stripPrefix     string
		stripSuffix     string
		replaceNewlines bool
	)

	cmd := &cobra.Command{
		Use:   "sanitize",
		Short: "Apply sanitization rules to secret values",
		Long:  "Sanitize reads secrets from the environment and applies strip/trim rules, printing the result.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecretSanitize(trimSpace, stripPrefix, stripSuffix, replaceNewlines)
		},
	}

	cmd.Flags().BoolVar(&trimSpace, "trim-space", false, "Trim leading and trailing whitespace from values")
	cmd.Flags().StringVar(&stripPrefix, "strip-prefix", "", "Strip a prefix from each value")
	cmd.Flags().StringVar(&stripSuffix, "strip-suffix", "", "Strip a suffix from each value")
	cmd.Flags().BoolVar(&replaceNewlines, "replace-newlines", false, "Replace newline characters with spaces")

	rootCmd.AddCommand(cmd)
}

func runSecretSanitize(trimSpace bool, stripPrefix, stripSuffix string, replaceNewlines bool) error {
	secrets := map[string]string{}
	for _, e := range os.Environ() {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				secrets[e[:i]] = e[i+1:]
				break
			}
		}
	}

	s := secretsanitize.New(secretsanitize.Rule{
		TrimSpace:       trimSpace,
		StripPrefix:     stripPrefix,
		StripSuffix:     stripSuffix,
		ReplaceNewlines: replaceNewlines,
	})

	out, err := s.Apply(secrets)
	if err != nil {
		return fmt.Errorf("sanitize: %w", err)
	}

	for k, v := range out {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
	}
	return nil
}
