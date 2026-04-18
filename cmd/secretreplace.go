package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secretreplace"
)

var secretreplaceRules []string

func init() {
	secretreplaceCmd := &cobra.Command{
		Use:   "secretreplace",
		Short: "Apply find-and-replace rules to secret values",
		Long: `Reads secrets from the environment and applies find-and-replace
rules to every value. Rules are specified as find=replace pairs.

Example:
  vaultpipe secretreplace --rule prod=dev --rule old=new`,
		RunE: runSecretreplace,
	}
	secretreplaceCmd.Flags().StringArrayVar(&secretreplaceRules, "rule", nil, "find=replace rule (repeatable)")
	_ = secretreplaceCmd.MarkFlagRequired("rule")
	rootCmd.AddCommand(secretreplaceCmd)
}

func runSecretreplace(cmd *cobra.Command, _ []string) error {
	rules, err := secretreplace.ParseRules(secretreplaceRules)
	if err != nil {
		return fmt.Errorf("secretreplace: %w", err)
	}

	r, err := secretreplace.New(rules)
	if err != nil {
		return fmt.Errorf("secretreplace: %w", err)
	}

	secrets := map[string]string{}
	for _, e := range os.Environ() {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				secrets[e[:i]] = e[i+1:]
				break
			}
		}
	}

	out, err := r.Apply(secrets)
	if err != nil {
		return fmt.Errorf("secretreplace: %w", err)
	}

	for k, v := range out {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
	}
	return nil
}
