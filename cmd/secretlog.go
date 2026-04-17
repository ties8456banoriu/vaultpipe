package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var secretlogCmd = &cobra.Command{
	Use:   "secretlog",
	Short: "Display the in-process secret access log (JSON)",
	Long: `Prints all secret access entries recorded during the current
vaultpipe session as a JSON array. Each entry includes the env key,
vault path, optional profile, and access timestamp.`,
	RunE: runSecretlog,
}

func init() {
	rootCmd.AddCommand(secretlogCmd)
	secretlogCmd.Flags().StringP("key", "k", "", "Filter entries by env key")
}

func runSecretlog(cmd *cobra.Command, _ []string) error {
	keyFilter, _ := cmd.Flags().GetString("key")

	if sharedSecretLog == nil {
		return fmt.Errorf("secret log is not initialised")
	}

	var entries interface{}
	if keyFilter != "" {
		entries = sharedSecretLog.ForKey(keyFilter)
	} else {
		entries = sharedSecretLog.All()
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
