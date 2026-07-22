package renderer

import (
	"fmt"
	"io"
	"sort"

	"github.com/jedi-knights/repo/internal/domain"
)

// Text is the human-readable renderer used as the default. Output is
// deliberately plain (no ANSI colors, no unicode boxes) so downstream
// pipes and CI logs handle it cleanly.
type Text struct{}

// NewText returns a Text renderer.
func NewText() *Text { return &Text{} }

// Render writes a plain-text view of report and summary to w.
func (Text) Render(w io.Writer, report *domain.Report, summary domain.Summary) error {
	if _, err := fmt.Fprintf(w, "Root: %s\n", report.Root); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Components: %d\n", summary.Total); err != nil {
		return err
	}
	if len(summary.ByEcosystem) > 0 {
		if _, err := fmt.Fprintln(w, "By ecosystem:"); err != nil {
			return err
		}
		for _, k := range sortedKeys(summary.ByEcosystem) {
			if _, err := fmt.Fprintf(w, "  %s: %d\n", k, summary.ByEcosystem[k]); err != nil {
				return err
			}
		}
	}
	if len(summary.ByWorkspace) > 0 {
		if _, err := fmt.Fprintln(w, "By workspace:"); err != nil {
			return err
		}
		for _, k := range sortedKeys(summary.ByWorkspace) {
			if _, err := fmt.Fprintf(w, "  %s: %d\n", k, summary.ByWorkspace[k]); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fprintln(w, "Details:"); err != nil {
		return err
	}
	for _, c := range report.Components {
		if _, err := fmt.Fprintf(w, "  - %s at %s (confidence %.2f)\n", c.Kind, c.Root, c.Confidence); err != nil {
			return err
		}
		for _, ws := range c.Workspaces {
			if _, err := fmt.Fprintf(w, "      workspace: %s (%d members)\n", ws.Kind, len(ws.Members)); err != nil {
				return err
			}
		}
	}
	return nil
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
