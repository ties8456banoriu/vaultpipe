package cmd

import (
	"fmt"
	"strconv"

	"github.com/elizabethwanjiku703/vaultpipe/internal/secretpin"
	"github.com/spf13/cobra"
)

var globalPinner = secretpin.New()

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Manage secret version pins",
	}

	pinCmd.AddCommand(&cobra.Command{
		Use:   "set <env-key> <version>",
		Short: "Pin a secret to a specific Vault version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("version must be an integer: %w", err)
			}
			if err := globalPinner.Pin(args[0], version); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "pinned %s to version %d\n", args[0], version)
			return nil
		},
	})

	pinCmd.AddCommand(&cobra.Command{
		Use:   "remove <env-key>",
		Short: "Remove a version pin for a secret",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := globalPinner.Unpin(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "unpinned %s\n", args[0])
			return nil
		},
	})

	pinCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all pinned secrets",
		RunE: func(cmd *cobra.Command, args []string) error {
			all := globalPinner.All()
			if len(all) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no pins set")
				return nil
			}
			for _, p := range all {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\tv%d\t(pinned at %s)\n",
					p.EnvKey, p.Version, p.PinnedAt.Format("2006-01-02T15:04:05Z"))
			}
			return nil
		},
	})

	rootCmd.AddCommand(pinCmd)
}
