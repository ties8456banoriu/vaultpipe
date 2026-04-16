package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/elizamthomas1994/vaultpipe/internal/healthcheck"
	"github.com/spf13/cobra"
)

var healthcheckTimeout time.Duration

func init() {
	hcCmd := &cobra.Command{
		Use:   "healthcheck",
		Short: "Check connectivity to the Vault server",
		RunE:  runHealthcheck,
	}
	hcCmd.Flags().DurationVar(&healthcheckTimeout, "timeout", 5*time.Second, "request timeout")
	rootCmd.AddCommand(hcCmd)
}

func runHealthcheck(cmd *cobra.Command, _ []string) error {
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "http://127.0.0.1:8200"
	}

	checker := healthcheck.New(vaultAddr, healthcheckTimeout)
	res := checker.Check(cmd.Context())

	if res.Err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck failed: %v\n", res.Err)
		return res.Err
	}

	status := "reachable"
	if !res.Reachable {
		status = "unreachable"
	}

	fmt.Printf("vault %s | status=%d | latency=%s | addr=%s\n",
		status, res.StatusCode, res.Latency.Round(time.Millisecond), vaultAddr)

	if !res.Reachable {
		os.Exit(1)
	}
	return nil
}
