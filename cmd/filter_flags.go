package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/vaultdiff/internal/filter"
)

// filterOpts holds the parsed filter flags for the root command.
var filterOpts filter.Options

// registerFilterFlags attaches key-filtering flags to the given command.
func registerFilterFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&filterOpts.Prefix,
		"prefix",
		"",
		"only include keys with this prefix (e.g. \"db/\")",
	)
	cmd.Flags().BoolVar(
		&filterOpts.OnlyChanged,
		"only-changed",
		false,
		"exclude unchanged keys from output",
	)
	cmd.Flags().StringSliceVar(
		&filterOpts.Keys,
		"keys",
		nil,
		"comma-separated explicit key allowlist",
	)
}
