package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secretencode"
)

var (
	encodeEncoding string
	encodeDecode   bool
)

func init() {
	encodeCmd := &cobra.Command{
		Use:   "secretencode",
		Short: "Encode or decode secret values using base64 or base64url",
		RunE:  runSecretencode,
	}
	encodeCmd.Flags().StringVar(&encodeEncoding, "encoding", "base64", "Encoding to use: base64, base64url")
	encodeCmd.Flags().BoolVar(&encodeDecode, "decode", false, "Decode values instead of encoding")
	rootCmd.AddCommand(encodeCmd)
}

func runSecretencode(cmd *cobra.Command, args []string) error {
	enc := secretencode.Encoding(encodeEncoding)
	e, err := secretencode.New(enc, encodeDecode)
	if err != nil {
		return fmt.Errorf("secretencode: %w", err)
	}

	// Read secrets from environment for demonstration.
	secrets := map[string]string{}
	for _, key := range args {
		val := os.Getenv(key)
		if val == "" {
			fmt.Fprintf(cmd.ErrOrStderr(), "warning: key %q not found in environment\n", key)
			continue
		}
		secrets[key] = val
	}

	if len(secrets) == 0 {
		return fmt.Errorf("secretencode: no secrets to process; provide env var names as arguments")
	}

	out, err := e.Apply(secrets)
	if err != nil {
		return fmt.Errorf("secretencode: %w", err)
	}

	for k, v := range out {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
	}
	return nil
}
