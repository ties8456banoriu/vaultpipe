package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/celtechstarter/vaultpipe/internal/secretpriority"
	"github.com/spf13/cobra"
)

func init() {
	var entries []string

	cmd := &cobra.Command{
		Use:   "secretpriority",
		Short: "Resolve secrets by priority level",
		Long: `Accepts key=value@priority triples and prints the winning value per key.

Example:
  vaultpipe secretpriority --entry DB_PASS=low@1 --entry DB_PASS=high@10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecretpriority(entries)
		},
	}

	cmd.Flags().StringArrayVar(&entries, "entry", nil, "key=value@priority (repeatable)")
	_ = cmd.MarkFlagRequired("entry")
	rootCmd.AddCommand(cmd)
}

func runSecretpriority(entries []string) error {
	m := secretpriority.New()
	for _, raw := range entries {
		at := strings.LastIndex(raw, "@")
		if at < 0 {
			return fmt.Errorf("invalid entry %q: expected key=value@priority", raw)
		}
		kv, prioStr := raw[:at], raw[at+1:]
		eq := strings.Index(kv, "=")
		if eq < 0 {
			return fmt.Errorf("invalid entry %q: expected key=value", raw)
		}
		key, value := kv[:eq], kv[eq+1:]
		p, err := strconv.Atoi(prioStr)
		if err != nil {
			return fmt.Errorf("invalid priority %q: %w", prioStr, err)
		}
		if err := m.Add(key, value, secretpriority.Level(p)); err != nil {
			return err
		}
	}

	out, err := m.ResolveAll()
	if err != nil {
		return err
	}
	for k, v := range out {
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}
