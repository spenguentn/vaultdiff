package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/audit"
	"github.com/yourusername/vaultdiff/internal/compare"
	"github.com/yourusername/vaultdiff/internal/config"
	"github.com/yourusername/vaultdiff/internal/filter"
	"github.com/yourusername/vaultdiff/internal/report"
	"github.com/yourusername/vaultdiff/internal/vault"
)

var diffCmd = &cobra.Command{
	Use:   "diff <left-env> <right-env> <secret-path>",
	Short: "Diff a secret path between two Vault environments",
	Long: `Compare a KV secret at a given path across two configured Vault environments.

Examples:
  vaultdiff diff staging production secret/app/config
  vaultdiff diff staging production secret/app/config --only-changed
  vaultdiff diff staging production secret/app/config --mask-secrets --output json`,
	Args: cobra.ExactArgs(3),
	RunE: runDiff,
}

func init() {
	registerFilterFlags(diffCmd)
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	leftName := args[0]
	rightName := args[1]
	secretPath := args[2]

	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	leftEnv, err := cfg.Environment(leftName)
	if err != nil {
		return fmt.Errorf("resolving left environment %q: %w", leftName, err)
	}

	rightEnv, err := cfg.Environment(rightName)
	if err != nil {
		return fmt.Errorf("resolving right environment %q: %w", rightName, err)
	}

	pair, err := vault.NewEnvPair(leftEnv, rightEnv)
	if err != nil {
		return fmt.Errorf("building environment pair: %w", err)
	}

	engine := compare.NewEngine(pair)
	results, err := engine.Run(cmd.Context(), secretPath)
	if err != nil {
		return fmt.Errorf("comparing secrets at %q: %w", secretPath, err)
	}

	opts := filter.OptionsFromFlags(cmd.Flags())
	results = filter.Apply(results, opts)

	maskSecrets, _ := cmd.Flags().GetBool("mask-secrets")
	outputFormat, _ := cmd.Flags().GetString("output")

	r := report.New(pair, secretPath, results)

	renderer := report.NewRenderer(report.RenderOptions{
		Format:      outputFormat,
		MaskSecrets: maskSecrets,
	})

	if err := renderer.Render(os.Stdout, r); err != nil {
		return fmt.Errorf("rendering report: %w", err)
	}

	if cfg.Audit.Enabled {
		session := audit.NewSession(cfg.Audit.User)
		entry := session.BuildEntry(r.Summary(), secretPath, leftName, rightName)

		logger, err := audit.NewLogger(cfg.Audit.LogPath, cfg.Audit.Format)
		if err != nil {
			return fmt.Errorf("initialising audit logger: %w", err)
		}
		if err := logger.Write(entry); err != nil {
			return fmt.Errorf("writing audit log: %w", err)
		}
	}

	return nil
}
