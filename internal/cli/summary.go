package cli

import (
	"fmt"
	"sort"

	"github.com/jedi-knights/repo/internal/domain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newSummaryCmd builds the `repo summary <path>` command. Runs a scan
// but renders only the derived Summary — a much shorter output aimed
// at humans skimming a directory of unrelated repos.
func newSummaryCmd(deps Deps, v *viper.Viper) *cobra.Command {
	return &cobra.Command{
		Use:   "summary <path>",
		Short: "Print a short summary of ecosystems and workspaces in a repo.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			analyzer := domain.NewAnalyzer(deps.Scanner)
			report, err := analyzer.Analyze(cmd.Context(), args[0], scanOptsFromViper(v))
			if err != nil {
				return err
			}
			s := domain.Summarize(report)
			out := cmd.OutOrStdout()
			if _, err := fmt.Fprintf(out, "%s\n  components: %d\n  monorepo: %v\n",
				report.Root, s.Total, s.HasMonorepo); err != nil {
				return err
			}
			for _, k := range sortedIntMapKeys(s.ByEcosystem) {
				if _, err := fmt.Fprintf(out, "  %s: %d\n", k, s.ByEcosystem[k]); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func sortedIntMapKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
