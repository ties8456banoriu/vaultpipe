package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secretarchive"
)

var globalArchive = secretarchive.New()

func init() {
	archiveCmd := &cobra.Command{
		Use:   "archive",
		Short: "Manage named secret archives",
	}

	storeCmd := &cobra.Command{
		Use:   "store <name>",
		Short: "Store current secrets under a named archive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runArchiveStore(args[0])
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Retrieve a named archive and print its secrets",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runArchiveGet(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all named archives",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runArchiveList()
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a named archive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return globalArchive.Delete(args[0])
		},
	}

	archiveCmd.AddCommand(storeCmd, getCmd, listCmd, deleteCmd)
	rootCmd.AddCommand(archiveCmd)
}

func runArchiveStore(name string) error {
	// Placeholder: in production this would pull from the active secret set.
	fmt.Fprintf(os.Stdout, "archived secrets under %q\n", name)
	return nil
}

func runArchiveGet(name string) error {
	e, err := globalArchive.Get(name)
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVALUE")
	for k, v := range e.Secrets {
		fmt.Fprintf(w, "%s\t%s\n", k, v)
	}
	return w.Flush()
}

func runArchiveList() error {
	entries := globalArchive.All()
	if len(entries) == 0 {
		fmt.Fprintln(os.Stdout, "no archives found")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tARCHIVED AT\tKEYS")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%d\n", e.Name, e.ArchivedAt.Format("2006-01-02 15:04:05"), len(e.Secrets))
	}
	return w.Flush()
}
