package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secretclone"
)

var (
	clonePrefix    string
	cloneUppercase bool
)

func init() {
	cloneCmd := &cobra.Command{
		Use:   "secretclone",
		Short: "Clone and transform a secrets map",
		Long:  "Prints a transformed clone of the provided secrets for inspection.",
		RunE:  runSecretclone,
	}
	cloneCmd.Flags().StringVar(&clonePrefix, "prefix", "", "Prefix to prepend to every key")
	cloneCmd.Flags().BoolVar(&cloneUppercase, "uppercase", false, "Uppercase all keys")
	rootCmd.AddCommand(cloneCmd)
}

func runSecretclone(cmd *cobra.Command, args []string) error {
	// Demo: clone a hardcoded map; in real usage this comes from vault fetch.
	src := map[string]string{
		"db_password": "s3cr3t",
		"api_key":     "abc123",
	}

	opts := []secretclone.Option{}
	if clonePrefix != "" {
		opts = append(opts, secretclone.WithPrefix(clonePrefix))
	}
	if cloneUppercase {
		opts = append(opts, secretclone.WithUppercase())
	}

	c := secretclone.New(opts...)
	out, err := c.Clone(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}

	for k, v := range out {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
	}
	return nil
}
