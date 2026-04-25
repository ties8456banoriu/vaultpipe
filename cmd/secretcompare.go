package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secretcompare"
)

func init() {
	var pathA, pathB string

	cmd := &cobra.Command{
		Use:   "secret-compare",
		Short: "Compare two sets of secrets side by side",
		Long:  "Reads two .env files and compares their key/value pairs, reporting matches, mismatches, and missing keys.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecretCompare(pathA, pathB)
		},
	}

	cmd.Flags().StringVar(&pathA, "a", "", "Path to the first .env file (required)")
	cmd.Flags().StringVar(&pathB, "b", "", "Path to the second .env file (required)")
	_ = cmd.MarkFlagRequired("a")
	_ = cmd.MarkFlagRequired("b")

	rootCmd.AddCommand(cmd)
}

func runSecretCompare(pathA, pathB string) error {
	a, err := loadEnvFile(pathA)
	if err != nil {
		return fmt.Errorf("reading --a: %w", err)
	}
	b, err := loadEnvFile(pathB)
	if err != nil {
		return fmt.Errorf("reading --b: %w", err)
	}

	c := secretcompare.New()
	results, err := c.Compare(a, b)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tSTATUS\tVALUE_A\tVALUE_B")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Key, r.Status, r.ValueA, r.ValueB)
	}
	_ = w.Flush()

	fmt.Println()
	fmt.Println(secretcompare.Summary(results))
	return nil
}

// loadEnvFile reads a simple KEY=VALUE .env file into a map.
func loadEnvFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string)
	for _, line := range splitLines(string(data)) {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		for i, ch := range line {
			if ch == '=' {
				out[line[:i]] = line[i+1:]
				break
			}
		}
	}
	return out, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, ch := range s {
		if ch == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
