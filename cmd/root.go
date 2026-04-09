package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/config"
)

var (
	cfgFile    string
	appConfig  *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "vaultdiff",
	Short: "Diff and audit changes between HashiCorp Vault secret versions",
	Long: `vaultdiff is a CLI tool for comparing Vault secret versions
across environments and producing structured audit logs of changes.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		appConfig, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		return nil
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "",
		"config file (default: $HOME/.vaultdiff.yaml or ./vaultdiff.yaml)",
	)
	rootCmd.PersistentFlags().String(
		"output", "text",
		"output format: text or json",
	)
	rootCmd.PersistentFlags().Bool(
		"mask-secrets", true,
		"mask secret values in output",
	)
}
