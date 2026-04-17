package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secretmeta"
)

var secretmetaCmd = &cobra.Command{
	Use:   "secretmeta",
	Short: "Display metadata for tracked secrets",
	RunE:  runSecretmeta,
}

func init() {
	rootCmd.AddCommand(secretmetaCmd)
}

func runSecretmeta(cmd *cobra.Command, _ []string) error {
	// In a real implementation the store would be populated during the run
	// command and persisted; here we demonstrate the output format.
	store := secretmeta.New()

	// Seed with an example when running standalone so the table is visible.
	_ = store.Record(secretmeta.Meta{
		EnvKey:    "DB_PASSWORD",
		VaultPath: "secret/data/db",
		Mount:     "secret",
		Version:   2,
		FetchedAt: time.Now().UTC(),
	})

	all := store.All()
	if len(all) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no secret metadata recorded")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ENV KEY\tVAULT PATH\tMOUNT\tVERSION\tFETCHED AT")
	for _, m := range all {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
			m.EnvKey,
			m.VaultPath,
			m.Mount,
			m.Version,
			m.FetchedAt.Format(time.RFC3339),
		)
	}
	return w.Flush()
}
