package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jedi-knights/repo/internal/domain"
)

type fakeScanner struct {
	report *domain.Report
	err    error
	seen   domain.ScanOptions
	root   string
}

func (f *fakeScanner) Scan(_ context.Context, root string, opts domain.ScanOptions) (*domain.Report, error) {
	f.root = root
	f.seen = opts
	return f.report, f.err
}

func TestAnalyzer_ForwardsRootAndOptions(t *testing.T) {
	fs := &fakeScanner{report: &domain.Report{Root: "/x"}}
	a := domain.NewAnalyzer(fs)

	opts := domain.ScanOptions{MaxDepth: 7, MaxDirs: 100, MaxFileSize: 4096}
	_, err := a.Analyze(context.Background(), "/x", opts)
	if err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}
	if fs.root != "/x" {
		t.Fatalf("root forwarded as %q, want /x", fs.root)
	}
	if fs.seen != opts {
		t.Fatalf("opts forwarded as %+v, want %+v", fs.seen, opts)
	}
}

func TestAnalyzer_PropagatesScannerError(t *testing.T) {
	want := errors.New("boom")
	fs := &fakeScanner{err: want}
	a := domain.NewAnalyzer(fs)

	_, err := a.Analyze(context.Background(), "/x", domain.ScanOptions{})
	if !errors.Is(err, want) {
		t.Fatalf("err = %v, want to wrap %v", err, want)
	}
}

func TestAnalyzer_RejectsEmptyRoot(t *testing.T) {
	a := domain.NewAnalyzer(&fakeScanner{})
	_, err := a.Analyze(context.Background(), "", domain.ScanOptions{})
	if err == nil {
		t.Fatalf("expected error for empty root")
	}
}
