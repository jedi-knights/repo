// Package renderer implements ports.Renderer for a handful of output
// formats. Renderers are Strategy implementations; NewFactory is the
// Factory that dispatches by name. Add a new format by registering it
// in NewFactory — callers unchanged.
package renderer

import (
	"fmt"
	"sort"

	"github.com/jedi-knights/repo/internal/ports"
)

// Factory selects a Renderer by name.
type Factory struct {
	strategies map[string]ports.Renderer
}

// NewFactory returns a Factory pre-populated with the built-in
// strategies. Registration order is not significant; Formats() returns
// names in sorted order for stable UX (help text, error messages).
func NewFactory() *Factory {
	return &Factory{
		strategies: map[string]ports.Renderer{
			"text":  NewText(),
			"json":  NewJSON(),
			"table": NewTable(),
		},
	}
}

// For returns the Renderer registered for name. Unknown names produce
// an error naming the available formats so the caller can correct.
func (f *Factory) For(name string) (ports.Renderer, error) {
	r, ok := f.strategies[name]
	if !ok {
		return nil, fmt.Errorf("renderer: unknown format %q (known: %v)", name, f.Formats())
	}
	return r, nil
}

// Formats returns the registered format names in sorted order.
func (f *Factory) Formats() []string {
	names := make([]string, 0, len(f.strategies))
	for k := range f.strategies {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
