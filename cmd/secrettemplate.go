package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/secrettemplate"
)

func init() {
	var tmplStr string

	cmd := &cobra.Command{
		Use:   "secrettemplate",
		Short: "Render a Go template using secrets from the environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecrettemplate(tmplStr)
		},
	}
	cmd.Flags().StringVarP(&tmplStr, "template", "t", "", "Go template string to render (required)")
	_ = cmd.MarkFlagRequired("template")
	rootCmd.AddCommand(cmd)
}

func runSecrettemplate(tmplStr string) error {
	secrets := make(map[string]string)
	for _, e := range os.Environ() {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				secrets[e[:i]] = e[i+1:]
				break
			}
		}
	}

	r := secrettemplate.New()
	out, err := r.Render(tmplStr, secrets)
	if err != nil {
		return fmt.Errorf("secrettemplate: %w", err)
	}
	fmt.Println(out)
	return nil
}
