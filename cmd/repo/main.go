// Command repo is a demo CLI over the repometa library.
//
// Wiring only. All logic lives in internal/domain (pure), internal/ports
// (interfaces), and internal/adapters (concrete integrations).
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jedi-knights/repo/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	root := cli.NewRootCmd(cli.NewDefaultDeps())
	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "repo:", err)
		os.Exit(1)
	}
}
