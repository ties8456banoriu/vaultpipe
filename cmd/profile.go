package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpipe/internal/profile"
)

func defaultProfilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".vaultpipe", "profiles.json")
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage named Vault secret profiles",
}

var profileSetCmd = &cobra.Command{
	Use:   "set <name>",
	Short: "Create or update a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("secret-path")
		envFile, _ := cmd.Flags().GetString("env-file")
		s := profile.NewStore(defaultProfilePath())
		_ = s.Load()
		if err := s.Set(profile.Profile{
			Name:       args[0],
			SecretPath: path,
			EnvFile:    envFile,
		}); err != nil {
			return err
		}
		return s.Save()
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := profile.NewStore(defaultProfilePath())
		if err := s.Load(); err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tSECRET PATH\tENV FILE")
		for _, p := range s.List() {
			fmt.Fprintf(w, "%s\t%s\t%s\n", p.Name, p.SecretPath, p.EnvFile)
		}
		return w.Flush()
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s := profile.NewStore(defaultProfilePath())
		if err := s.Load(); err != nil {
			return err
		}
		if err := s.Delete(args[0]); err != nil {
			return err
		}
		return s.Save()
	},
}

func init() {
	profileSetCmd.Flags().StringP("secret-path", "p", "", "Vault secret path (required)")
	profileSetCmd.Flags().StringP("env-file", "e", ".env", "Output .env file path")
	_ = profileSetCmd.MarkFlagRequired("secret-path")

	profileCmd.AddCommand(profileSetCmd, profileListCmd, profileDeleteCmd)
	rootCmd.AddCommand(profileCmd)
}
