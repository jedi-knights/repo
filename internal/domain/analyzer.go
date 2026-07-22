package domain

import (
	"context"
	"errors"
	"fmt"
)

// ScanOptions is the shape of tunables the analyzer forwards to a
// Scanner. Zero values mean "use scanner default".
type ScanOptions struct {
	MaxDepth    int
	MaxDirs     int
	MaxFileSize int64
}

// Scanner performs the actual repository walk. Defined here (at the
// consumer) per Go convention; adapter packages implement it.
type Scanner interface {
	Scan(ctx context.Context, root string, opts ScanOptions) (*Report, error)
}

// Analyzer is the domain service that orchestrates a scan and derives
// downstream views. Holds a Scanner via constructor injection so tests
// can substitute a fake.
type Analyzer struct {
	scanner Scanner
}

// NewAnalyzer wires an Analyzer to a Scanner. Panics on nil scanner —
// programmer error, not a runtime condition, so a loud failure at
// startup is preferred over deferred nil-deref.
func NewAnalyzer(scanner Scanner) *Analyzer {
	if scanner == nil {
		panic("domain: NewAnalyzer requires a non-nil Scanner")
	}
	return &Analyzer{scanner: scanner}
}

// Analyze runs a scan of root with opts and returns the resulting
// Report. Errors from the scanner are wrapped for provenance.
func (a *Analyzer) Analyze(ctx context.Context, root string, opts ScanOptions) (*Report, error) {
	if root == "" {
		return nil, errors.New("domain: root path is empty")
	}
	rep, err := a.scanner.Scan(ctx, root, opts)
	if err != nil {
		return nil, fmt.Errorf("domain: scan %q: %w", root, err)
	}
	return rep, nil
}
