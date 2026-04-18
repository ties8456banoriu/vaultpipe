package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/vaultpipe/internal/secretrotate"
)

var globalRotator = secretrotate.New()

func init() {
	recordCmd := &cobra.Command{
		Use:   "rotate-record",
		Short: "Record a secret rotation event",
		RunE:  runRotateRecord,
	}
	recordCmd.Flags().String("key", "", "Env key (required)")
	recordCmd.Flags().String("path", "", "Vault path (required)")
	recordCmd.Flags().Int("version", 1, "Secret version")
	recordCmd.Flags().String("policy", "manual", "Rotation policy: manual or scheduled")

	listCmd := &cobra.Command{
		Use:   "rotate-list",
		Short: "List all rotation records",
		RunE:  runRotateList,
	}

	rootCmd.AddCommand(recordCmd, listCmd)
}

func runRotateRecord(cmd *cobra.Command, _ []string) error {
	key, _ := cmd.Flags().GetString("key")
	path, _ := cmd.Flags().GetString("path")
	version, _ := cmd.Flags().GetInt("version")
	policyStr, _ := cmd.Flags().GetString("policy")

	policy := secretrotate.PolicyManual
	if policyStr == "scheduled" {
		policy = secretrotate.PolicyScheduled
	}

	if err := globalRotator.Record(key, path, version, policy); err != nil {
		return fmt.Errorf("rotate-record: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "rotation recorded for %s (version %d)\n", key, version)
	return nil
}

func runRotateList(cmd *cobra.Command, _ []string) error {
	all := globalRotator.All()
	if len(all) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no rotation records found")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVAULT PATH\tVERSION\tPOLICY\tROTATED AT")
	for _, recs := range all {
		for _, rec := range recs {
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n",
				rec.EnvKey, rec.VaultPath, rec.Version,
				rec.Policy, rec.RotatedAt.Format(time.RFC3339))
		}
	}
	return w.Flush()
}
