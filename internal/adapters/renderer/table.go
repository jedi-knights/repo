package renderer

import (
	"fmt"
	"io"
	"strings"

	"github.com/jedi-knights/repo/internal/domain"
	"github.com/olekukonko/tablewriter"
)

// Table renders the component list as an ASCII table. Useful in
// terminals and code review comments; not intended for machine
// consumption.
type Table struct{}

// NewTable returns a Table renderer.
func NewTable() *Table { return &Table{} }

// Render writes an ASCII table of components to w. The summary is
// emitted as a small header block above the table so scanning by eye
// mirrors the JSON payload.
func (Table) Render(w io.Writer, report *domain.Report, summary domain.Summary) error {
	if _, err := fmt.Fprintf(w, "Root: %s   Components: %d   Monorepo: %v\n\n",
		report.Root, summary.Total, summary.HasMonorepo); err != nil {
		return err
	}

	t := tablewriter.NewTable(w)
	t.Header("Kind", "Root", "Confidence", "Workspaces")
	for _, c := range report.Components {
		if err := t.Append(c.Kind, c.Root, fmt.Sprintf("%.2f", c.Confidence), workspaceCell(c.Workspaces)); err != nil {
			return err
		}
	}
	return t.Render()
}

func workspaceCell(ws []domain.Workspace) string {
	if len(ws) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(ws))
	for _, w := range ws {
		parts = append(parts, w.Kind)
	}
	return strings.Join(parts, ", ")
}
