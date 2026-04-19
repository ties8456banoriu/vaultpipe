package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secrettrim"
)

func init() {
	var rulesFlag []string

	cmd := &cobra.Command{
		Use:   "secrettrim",
		Short: "Trim secret values by start/end index",
		Long: `Apply substring trimming rules to secret values.

Each rule is specified as KEY:START:END where END may be -1 for end-of-string.

Example:
  vaultpipe secrettrim --rule TOKEN:7:-1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecrettrim(rulesFlag)
		},
	}

	cmd.Flags().StringArrayVar(&rulesFlag, "rule", nil, "Trim rule as KEY:START:END (END=-1 for full tail)")
	_ = cmd.MarkFlagRequired("rule")

	rootCmd.AddCommand(cmd)
}

func runSecrettrim(rawRules []string) error {
	rules, err := parseSecrettrimRules(rawRules)
	if err != nil {
		return fmt.Errorf("secrettrim: %w", err)
	}
	tr, err := secrettrim.New(rules)
	if err != nil {
		return fmt.Errorf("secrettrim: %w", err)
	}
	// Placeholder: in a real pipeline this would receive secrets from context.
	_ = tr
	fmt.Println("secrettrim: rules registered successfully")
	return nil
}

func parseSecrettrimRules(raw []string) ([]secrettrim.Rule, error) {
	var rules []secrettrim.Rule
	for _, r := range raw {
		parts := strings.SplitN(r, ":", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid rule %q: expected KEY:START:END", r)
		}
		start, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid start in rule %q: %w", r, err)
		}
		end, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid end in rule %q: %w", r, err)
		}
		rules = append(rules, secrettrim.Rule{Key: parts[0], Start: start, End: end})
	}
	return rules, nil
}
