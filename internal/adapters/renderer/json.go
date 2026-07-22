package renderer

import (
	"encoding/json"
	"io"

	"github.com/jedi-knights/repo/internal/domain"
)

// JSON emits report + summary as an indented JSON object. Consumers
// that want a machine-readable feed target this format.
type JSON struct{}

// NewJSON returns a JSON renderer.
func NewJSON() *JSON { return &JSON{} }

type jsonPayload struct {
	Root       string             `json:"root"`
	Stats      domain.Stats       `json:"stats"`
	Components []domain.Component `json:"components"`
	Summary    domain.Summary     `json:"summary"`
}

// Render writes the JSON payload to w.
func (JSON) Render(w io.Writer, report *domain.Report, summary domain.Summary) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jsonPayload{
		Root:       report.Root,
		Stats:      report.Stats,
		Components: report.Components,
		Summary:    summary,
	})
}
