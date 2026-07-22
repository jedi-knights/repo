package renderer_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/jedi-knights/repo/internal/adapters/renderer"
	"github.com/jedi-knights/repo/internal/domain"
)

func sampleReport() (*domain.Report, domain.Summary) {
	r := &domain.Report{
		Root: "/tmp/x",
		Components: []domain.Component{
			{Kind: "go-module", Root: ".", Confidence: 1.0},
			{Kind: "node-package", Root: "web", Confidence: 1.0, Workspaces: []domain.Workspace{{Kind: "pnpm-workspace"}}},
		},
		Stats: domain.Stats{DirsVisited: 12, FilesSeen: 40},
	}
	return r, domain.Summarize(r)
}

func TestFactory_KnownFormats(t *testing.T) {
	f := renderer.NewFactory()
	for _, name := range f.Formats() {
		if _, err := f.For(name); err != nil {
			t.Errorf("Factory.For(%q) error: %v", name, err)
		}
	}
}

func TestFactory_UnknownFormatErrors(t *testing.T) {
	f := renderer.NewFactory()
	if _, err := f.For("xml"); err == nil {
		t.Fatalf("expected error for unknown format")
	}
}

func TestTextRenderer_ContainsRootAndComponents(t *testing.T) {
	r, s := sampleReport()
	var buf bytes.Buffer
	if err := renderer.NewText().Render(&buf, r, s); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"/tmp/x", "go-module", "node-package", "pnpm-workspace"} {
		if !strings.Contains(out, want) {
			t.Errorf("text output missing %q: %s", want, out)
		}
	}
}

func TestJSONRenderer_ParsesBackWithExpectedFields(t *testing.T) {
	r, s := sampleReport()
	var buf bytes.Buffer
	if err := renderer.NewJSON().Render(&buf, r, s); err != nil {
		t.Fatalf("Render: %v", err)
	}
	var got struct {
		Root       string             `json:"root"`
		Components []domain.Component `json:"components"`
		Summary    domain.Summary     `json:"summary"`
	}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v; raw=%s", err, buf.String())
	}
	if got.Root != "/tmp/x" {
		t.Errorf("root = %q, want /tmp/x", got.Root)
	}
	if len(got.Components) != 2 {
		t.Errorf("components = %d, want 2", len(got.Components))
	}
	if got.Summary.Total != 2 {
		t.Errorf("summary.total = %d, want 2", got.Summary.Total)
	}
}

func TestTableRenderer_HasHeadersAndRows(t *testing.T) {
	r, s := sampleReport()
	var buf bytes.Buffer
	if err := renderer.NewTable().Render(&buf, r, s); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()
	// Case-insensitive check — tablewriter may upper- or title-case headers.
	lower := strings.ToLower(out)
	for _, want := range []string{"kind", "root", "go-module", "node-package"} {
		if !strings.Contains(lower, want) {
			t.Errorf("table output missing %q: %s", want, out)
		}
	}
}
