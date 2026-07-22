// Package cli builds the cobra command tree and wires domain and
// adapter dependencies. Configuration is layered by viper: defaults,
// config file, then env vars (REPO_*). Flags override the layers below.
package cli

import (
	"fmt"
	"strings"

	"github.com/jedi-knights/repo/internal/adapters/renderer"
	"github.com/jedi-knights/repo/internal/adapters/scanner"
	"github.com/jedi-knights/repo/internal/domain"
	"github.com/jedi-knights/repo/internal/ports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Deps groups the pluggable pieces the CLI needs. Constructed by
// NewDefaultDeps for production wiring; tests supply their own to
// exercise handlers without hitting the filesystem.
type Deps struct {
	Scanner  domain.Scanner
	Renderer ports.RendererFactory
}

// NewDefaultDeps returns Deps wired to the production adapters.
func NewDefaultDeps() Deps {
	return Deps{
		Scanner:  scanner.NewRepometa(),
		Renderer: renderer.NewFactory(),
	}
}

// NewRootCmd assembles the full command tree bound to deps. Extracted
// as a function (rather than a package-level `rootCmd`) so tests can
// build a fresh tree per case and inspect its output.
func NewRootCmd(deps Deps) *cobra.Command {
	v := viper.New()
	v.SetEnvPrefix("REPO")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	root := &cobra.Command{
		Use:           "repo",
		Short:         "Analyze arbitrary source repositories via repometa.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().String("config", "", "path to a viper config file (yaml/toml/json)")
	root.PersistentFlags().String("format", "text",
		fmt.Sprintf("output format (%s)", strings.Join(deps.Renderer.Formats(), ", ")))
	root.PersistentFlags().Int("max-depth", 0, "override scanner max directory depth (0 = library default)")
	root.PersistentFlags().Int("max-dirs", 0, "override scanner max directory count (0 = library default)")
	root.PersistentFlags().Int64("max-file-size", 0, "override scanner max file size in bytes (0 = library default)")

	// Bind flags so viper sees them alongside env / file layers.
	_ = v.BindPFlag("format", root.PersistentFlags().Lookup("format"))
	_ = v.BindPFlag("max-depth", root.PersistentFlags().Lookup("max-depth"))
	_ = v.BindPFlag("max-dirs", root.PersistentFlags().Lookup("max-dirs"))
	_ = v.BindPFlag("max-file-size", root.PersistentFlags().Lookup("max-file-size"))

	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		cfg, _ := cmd.Flags().GetString("config")
		if cfg == "" {
			return nil
		}
		v.SetConfigFile(cfg)
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("read config %q: %w", cfg, err)
		}
		return nil
	}

	root.AddCommand(newScanCmd(deps, v))
	root.AddCommand(newSummaryCmd(deps, v))
	root.AddCommand(newVersionCmd())
	return root
}

// scanOptsFromViper reads the bounded-traversal knobs from v.
// Zero values are preserved so the scanner adapter forwards
// "use library default" downstream.
func scanOptsFromViper(v *viper.Viper) domain.ScanOptions {
	return domain.ScanOptions{
		MaxDepth:    v.GetInt("max-depth"),
		MaxDirs:     v.GetInt("max-dirs"),
		MaxFileSize: v.GetInt64("max-file-size"),
	}
}
