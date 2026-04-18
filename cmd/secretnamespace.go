package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secretnamespace"
)

var nsManager = secretnamespace.New()

func init() {
	nsSetCmd := &cobra.Command{
		Use:   "ns-set <namespace> <KEY=VALUE>...",
		Short: "Store secrets under a namespace",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runNsSet,
	}
	nsGetCmd := &cobra.Command{
		Use:   "ns-get <namespace>",
		Short: "Retrieve secrets from a namespace",
		Args:  cobra.ExactArgs(1),
		RunE:  runNsGet,
	}
	nsListCmd := &cobra.Command{
		Use:   "ns-list",
		Short: "List all registered namespaces",
		RunE:  runNsList,
	}
	nsMergeCmd := &cobra.Command{
		Use:   "ns-merge <namespace>...",
		Short: "Merge secrets from multiple namespaces",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runNsMerge,
	}
	rootCmd.AddCommand(nsSetCmd, nsGetCmd, nsListCmd, nsMergeCmd)
}

func runNsSet(cmd *cobra.Command, args []string) error {
	ns := args[0]
	secrets := make(map[string]string)
	for _, pair := range args[1:] {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid key=value pair: %s", pair)
		}
		secrets[parts[0]] = parts[1]
	}
	if err := nsManager.Set(ns, secrets); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "namespace %q set with %d key(s)\n", ns, len(secrets))
	return nil
}

func runNsGet(cmd *cobra.Command, args []string) error {
	secrets, err := nsManager.Get(args[0])
	if err != nil {
		return err
	}
	for k, v := range secrets {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
	}
	return nil
}

func runNsList(cmd *cobra.Command, args []string) error {
	list := nsManager.List()
	if len(list) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no namespaces registered")
		return nil
	}
	for _, ns := range list {
		fmt.Fprintln(cmd.OutOrStdout(), ns)
	}
	return nil
}

func runNsMerge(cmd *cobra.Command, args []string) error {
	result, err := nsManager.Merge(args...)
	if err != nil {
		return err
	}
	for k, v := range result {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
	}
	return nil
}
