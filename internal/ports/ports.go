// Package ports declares interfaces that the CLI layer needs from
// pluggable strategies. Interfaces consumed by the domain live in the
// domain package itself (Go convention: define at the consumer).
package ports

import (
	"io"

	"github.com/jedi-knights/repo/internal/domain"
)

// Renderer serializes a Report (and its Summary) to w. Each concrete
// renderer is a Strategy; RendererFactory selects among them.
type Renderer interface {
	Render(w io.Writer, report *domain.Report, summary domain.Summary) error
}

// RendererFactory returns a Renderer for a named format, or an error
// when the format is not registered.
type RendererFactory interface {
	For(format string) (Renderer, error)
	Formats() []string
}
