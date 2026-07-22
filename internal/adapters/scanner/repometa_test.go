package scanner_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jedi-knights/repo/internal/adapters/scanner"
	"github.com/jedi-knights/repo/internal/domain"
)

// makeGoModuleDir creates a minimal Go module at dir so repometa detects it.
func makeGoModuleDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module x\n\ngo 1.26\n"), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	return dir
}

func TestRepometa_DetectsGoModule(t *testing.T) {
	root := makeGoModuleDir(t)
	s := scanner.NewRepometa()

	rep, err := s.Scan(context.Background(), root, domain.ScanOptions{})
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(rep.Components) != 1 {
		t.Fatalf("component count = %d, want 1: %+v", len(rep.Components), rep.Components)
	}
	c := rep.Components[0]
	if c.Kind != "go-module" {
		t.Fatalf("Kind = %q, want go-module", c.Kind)
	}
}

func TestRepometa_RejectsEmptyRoot(t *testing.T) {
	s := scanner.NewRepometa()
	if _, err := s.Scan(context.Background(), "", domain.ScanOptions{}); err == nil {
		t.Fatalf("expected error for empty root")
	}
}

func TestRepometa_ForwardsScanOptions(t *testing.T) {
	root := makeGoModuleDir(t)
	// Nested dir so MaxDepth=0 would still catch root, but depth cap
	// exercises the option path; correctness of the cap itself is
	// repometa's contract, not ours.
	if err := os.MkdirAll(filepath.Join(root, "a", "b", "c"), 0o755); err != nil {
		t.Fatal(err)
	}
	s := scanner.NewRepometa()
	_, err := s.Scan(context.Background(), root, domain.ScanOptions{
		MaxDepth: 1, MaxDirs: 10, MaxFileSize: 1024,
	})
	if err != nil {
		t.Fatalf("Scan with opts: %v", err)
	}
}
