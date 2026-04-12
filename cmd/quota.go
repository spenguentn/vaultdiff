package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var (
	quotaMount      string
	quotaPrefix     string
	quotaMaxReads   int
	quotaWindowSecs int
)

var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Manage secret read quotas",
	Long:  "Register and inspect rate-limit quotas applied to Vault secret paths.",
}

var quotaRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new secret quota",
	RunE:  runQuotaRegister,
}

func init() {
	quotaRegisterCmd.Flags().StringVar(&quotaMount, "mount", "", "Vault mount path (required)")
	quotaRegisterCmd.Flags().StringVar(&quotaPrefix, "prefix", "", "Secret path prefix")
	quotaRegisterCmd.Flags().IntVar(&quotaMaxReads, "max-reads", 100, "Maximum reads per window")
	quotaRegisterCmd.Flags().IntVar(&quotaWindowSecs, "window", 60, "Window size in seconds")
	_ = quotaRegisterCmd.MarkFlagRequired("mount")

	quotaCmd.AddCommand(quotaRegisterCmd)
	rootCmd.AddCommand(quotaCmd)
}

func runQuotaRegister(cmd *cobra.Command, _ []string) error {
	q := vault.SecretQuota{
		Mount:      quotaMount,
		Prefix:     quotaPrefix,
		Scope:      vault.QuotaScopeMount,
		MaxReads:   quotaMaxReads,
		WindowSize: time.Duration(quotaWindowSecs) * time.Second,
	}
	if err := q.Validate(); err != nil {
		return fmt.Errorf("invalid quota: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(),
		"quota registered: mount=%s prefix=%s max_reads=%d window=%ds\n",
		q.Mount, q.Prefix, q.MaxReads, quotaWindowSecs,
	)
	return nil
}
