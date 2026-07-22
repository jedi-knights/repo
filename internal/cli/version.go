package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the string emitted by `repo version`. Overridden at
// link time via -ldflags "-X github.com/jedi-knights/repo/internal/cli.Version=v0.1.0".
var Version = "dev"

// Commit is the short git SHA of the build, when set via -ldflags.
var Commit = "none"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the repo binary version.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "repo %s (%s)\n", Version, Commit)
			return err
		},
	}
}
