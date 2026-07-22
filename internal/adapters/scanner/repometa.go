// Package scanner adapts the repometa library to the domain.Scanner
// port. This is the only package in the tree that imports repometa;
// swapping the analysis backend is confined to a single file.
package scanner

import (
	"context"

	"github.com/jedi-knights/repo/internal/domain"
	"github.com/jedi-knights/repometa"
)

// Repometa is the adapter that satisfies domain.Scanner via the
// upstream repometa library.
type Repometa struct{}

// NewRepometa constructs a Repometa scanner. Kept as a constructor
// (rather than exposing the zero value) to leave room for future
// dependencies (metrics client, cache, etc.) without a breaking change.
func NewRepometa() *Repometa { return &Repometa{} }

// Scan runs repometa.Scan against root and translates the result into
// the domain shape. The ctx is not forwarded because repometa.Scan is
// synchronous and un-cancellable in its current version; callers that
// need timeouts should wrap the call in a goroutine + select.
func (Repometa) Scan(_ context.Context, root string, opts domain.ScanOptions) (*domain.Report, error) {
	m, err := repometa.Scan(root, toRepometaOptions(opts)...)
	if err != nil {
		return nil, err
	}
	return fromManifest(m), nil
}

// toRepometaOptions turns zero-valued options into "use library default"
// by simply omitting them. A zero cap otherwise makes the scan return
// nothing useful.
func toRepometaOptions(o domain.ScanOptions) []repometa.Option {
	var opts []repometa.Option
	if o.MaxDepth > 0 {
		opts = append(opts, repometa.WithMaxDepth(o.MaxDepth))
	}
	if o.MaxDirs > 0 {
		opts = append(opts, repometa.WithMaxDirs(o.MaxDirs))
	}
	if o.MaxFileSize > 0 {
		opts = append(opts, repometa.WithMaxFileSize(o.MaxFileSize))
	}
	return opts
}

func fromManifest(m *repometa.Manifest) *domain.Report {
	if m == nil {
		return &domain.Report{}
	}
	comps := make([]domain.Component, 0, len(m.Components))
	for _, c := range m.Components {
		comps = append(comps, fromComponent(c))
	}
	return &domain.Report{
		Root:       m.Root,
		Components: comps,
		Stats: domain.Stats{
			DirsVisited:     m.Stats.DirsVisited,
			FilesSeen:       m.Stats.FilesSeen,
			DepthCapHits:    m.Stats.DepthCapHits,
			DirCapHits:      m.Stats.DirCapHits,
			SymlinksSkipped: m.Stats.SymlinksSkipped,
		},
	}
}

func fromComponent(c repometa.Component) domain.Component {
	ev := make([]domain.Evidence, 0, len(c.Evidence))
	for _, e := range c.Evidence {
		ev = append(ev, domain.Evidence{Path: e.Path, Reason: e.Reason})
	}
	ws := make([]domain.Workspace, 0, len(c.Workspaces))
	for _, w := range c.Workspaces {
		ws = append(ws, domain.Workspace{Kind: string(w.Kind), Members: append([]string(nil), w.Members...)})
	}
	return domain.Component{
		Kind:       string(c.Kind),
		Root:       c.Root,
		Evidence:   ev,
		Confidence: c.Confidence,
		Workspaces: ws,
		Attributes: c.Attributes,
	}
}
