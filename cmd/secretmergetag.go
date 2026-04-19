package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/vaultpipe/internal/secretmergetag"
)

func init() {
	var policy string
	var files []string

	cmd := &cobra.Command{
		Use:   "secretmergetag",
		Short: "Merge tag maps from multiple JSON files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecretmergetag(files, secretmergetag.ConflictPolicy(policy))
		},
	}
	cmd.Flags().StringVar(&policy, "policy", "skip", "Conflict policy: skip, overwrite, error")
	cmd.Flags().StringArrayVar(&files, "file", nil, "JSON tag map files to merge (repeatable)")
	_ = cmd.MarkFlagRequired("file")
	rootCmd.AddCommand(cmd)
}

func runSecretmergetag(files []string, policy secretmergetag.ConflictPolicy) error {
	merger, err := secretmergetag.New(policy)
	if err != nil {
		return err
	}
	var sources []map[string]map[string]string
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("secretmergetag: reading %q: %w", f, err)
		}
		var m map[string]map[string]string
		if err := json.Unmarshal(data, &m); err != nil {
			return fmt.Errorf("secretmergetag: parsing %q: %w", f, err)
		}
		sources = append(sources, m)
	}
	out, err := merger.Merge(sources...)
	if err != nil {
		return err
	}
	return json.NewEncoder(os.Stdout).Encode(out)
}
