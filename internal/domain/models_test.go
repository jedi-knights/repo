package domain_test

import (
	"testing"

	"github.com/jedi-knights/repo/internal/domain"
)

func TestSummarize_EmptyReport(t *testing.T) {
	r := &domain.Report{Root: "/tmp/x"}
	s := domain.Summarize(r)
	if s.Total != 0 {
		t.Fatalf("Total = %d, want 0", s.Total)
	}
	if s.HasMonorepo {
		t.Fatalf("HasMonorepo = true, want false")
	}
	if len(s.ByEcosystem) != 0 || len(s.ByWorkspace) != 0 {
		t.Fatalf("expected empty maps, got %+v %+v", s.ByEcosystem, s.ByWorkspace)
	}
}

func TestSummarize_CountsByEcosystemAndWorkspace(t *testing.T) {
	r := &domain.Report{
		Root: "/tmp/x",
		Components: []domain.Component{
			{Kind: "go-module", Root: "."},
			{Kind: "go-module", Root: "svc/api"},
			{Kind: "node-package", Root: "web", Workspaces: []domain.Workspace{
				{Kind: "pnpm-workspace", Members: []string{"web/apps/a"}},
			}},
		},
	}
	s := domain.Summarize(r)
	if s.Total != 3 {
		t.Fatalf("Total = %d, want 3", s.Total)
	}
	if s.ByEcosystem["go-module"] != 2 {
		t.Fatalf("go-module count = %d, want 2", s.ByEcosystem["go-module"])
	}
	if s.ByEcosystem["node-package"] != 1 {
		t.Fatalf("node-package count = %d, want 1", s.ByEcosystem["node-package"])
	}
	if s.ByWorkspace["pnpm-workspace"] != 1 {
		t.Fatalf("pnpm-workspace count = %d, want 1", s.ByWorkspace["pnpm-workspace"])
	}
	if !s.HasMonorepo {
		t.Fatalf("HasMonorepo = false, want true")
	}
}

func TestSummarize_MultipleWorkspacesOnOneComponent(t *testing.T) {
	r := &domain.Report{
		Components: []domain.Component{
			{Kind: "node-package", Workspaces: []domain.Workspace{
				{Kind: "pnpm-workspace"},
				{Kind: "turborepo"},
			}},
		},
	}
	s := domain.Summarize(r)
	if s.ByWorkspace["pnpm-workspace"] != 1 || s.ByWorkspace["turborepo"] != 1 {
		t.Fatalf("expected both workspaces counted, got %+v", s.ByWorkspace)
	}
}
