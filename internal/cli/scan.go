package cli

import (
	"fmt"

	"github.com/jedi-knights/repo/internal/domain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newScanCmd builds the `repo scan <path>` command. Runs a full scan
// and renders the report via the format selected on the root command.
func newScanCmd(deps Deps, v *viper.Viper) *cobra.Command {
	return &cobra.Command{
		Use:   "scan <path>",
		Short: "Scan a repository and render the detected components.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := args[0]
			analyzer := domain.NewAnalyzer(deps.Scanner)
			report, err := analyzer.Analyze(cmd.Context(), root, scanOptsFromViper(v))
			if err != nil {
				return err
			}
			format := v.GetString("format")
			r, err := deps.Renderer.For(format)
			if err != nil {
				return err
			}
			if err := r.Render(cmd.OutOrStdout(), report, domain.Summarize(report)); err != nil {
				return fmt.Errorf("render %s: %w", format, err)
			}
			return nil
		},
	}
}
