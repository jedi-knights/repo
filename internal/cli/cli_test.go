package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/jedi-knights/repo/internal/adapters/renderer"
	"github.com/jedi-knights/repo/internal/cli"
	"github.com/jedi-knights/repo/internal/domain"
)

type stubScanner struct {
	rep *domain.Report
	err error
}

func (s *stubScanner) Scan(_ context.Context, _ string, _ domain.ScanOptions) (*domain.Report, error) {
	return s.rep, s.err
}

func testDeps(rep *domain.Report, err error) cli.Deps {
	return cli.Deps{
		Scanner:  &stubScanner{rep: rep, err: err},
		Renderer: renderer.NewFactory(),
	}
}

func runCmd(t *testing.T, deps cli.Deps, args ...string) (string, string, error) {
	t.Helper()
	root := cli.NewRootCmd(deps)
	var out, errBuf bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errBuf)
	root.SetArgs(args)
	err := root.Execute()
	return out.String(), errBuf.String(), err
}

func sampleReport() *domain.Report {
	return &domain.Report{
		Root: "/tmp/x",
		Components: []domain.Component{
			{Kind: "go-module", Root: ".", Confidence: 1.0},
		},
	}
}

func TestScan_TextFormatDefault(t *testing.T) {
	out, _, err := runCmd(t, testDeps(sampleReport(), nil), "scan", "/tmp/x")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(out, "go-module") {
		t.Errorf("expected text output to mention go-module, got: %s", out)
	}
}

func TestScan_JSONFormat(t *testing.T) {
	out, _, err := runCmd(t, testDeps(sampleReport(), nil), "scan", "--format", "json", "/tmp/x")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	var payload struct {
		Root string `json:"root"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("output not JSON: %v; raw=%s", err, out)
	}
	if payload.Root != "/tmp/x" {
		t.Errorf("root = %q, want /tmp/x", payload.Root)
	}
}

func TestScan_UnknownFormatReturnsError(t *testing.T) {
	_, _, err := runCmd(t, testDeps(sampleReport(), nil), "scan", "--format", "xml", "/tmp/x")
	if err == nil {
		t.Fatalf("expected error for unknown format")
	}
	if !strings.Contains(err.Error(), "unknown format") {
		t.Errorf("error = %v, want mention of unknown format", err)
	}
}

func TestScan_ScannerErrorPropagates(t *testing.T) {
	want := errors.New("scanner down")
	_, _, err := runCmd(t, testDeps(nil, want), "scan", "/tmp/x")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "scanner down") {
		t.Errorf("error = %v, want to wrap %v", err, want)
	}
}

func TestSummary_PrintsRootAndCounts(t *testing.T) {
	out, _, err := runCmd(t, testDeps(sampleReport(), nil), "summary", "/tmp/x")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	for _, want := range []string{"/tmp/x", "components: 1", "go-module: 1"} {
		if !strings.Contains(out, want) {
			t.Errorf("summary missing %q: %s", want, out)
		}
	}
}

func TestVersion_PrintsBuildInfo(t *testing.T) {
	out, _, err := runCmd(t, testDeps(nil, nil), "version")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.HasPrefix(out, "repo ") {
		t.Errorf("version output = %q, want prefix 'repo '", out)
	}
}
