// Package domain holds pure business types for the repo analyzer.
//
// Nothing here imports adapters or the CLI layer. Types deliberately
// mirror the shape of repometa.Manifest but stay independent of it so
// that a future scanner backend can be swapped in without touching
// domain code.
package domain

// Report is the domain view of a scanned repository.
type Report struct {
	Root       string      `json:"root"`
	Components []Component `json:"components"`
	Stats      Stats       `json:"stats"`
}

// Component is a single buildable unit anchored at a directory inside
// Report.Root. Root paths are slash-separated and use "." for the repo
// root, matching the underlying scanner's convention.
type Component struct {
	Kind       string            `json:"kind"`
	Root       string            `json:"root"`
	Evidence   []Evidence        `json:"evidence,omitempty"`
	Confidence float64           `json:"confidence"`
	Workspaces []Workspace       `json:"workspaces,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Evidence records one fact that led to a Component being reported.
type Evidence struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

// Workspace describes a monorepo layout anchored at a Component.
type Workspace struct {
	Kind    string   `json:"kind"`
	Members []string `json:"members,omitempty"`
}

// Stats mirrors the scanner's traversal counters so callers can tell
// when a cap was hit and the scan should be widened.
type Stats struct {
	DirsVisited     int `json:"dirs_visited"`
	FilesSeen       int `json:"files_seen"`
	DepthCapHits    int `json:"depth_cap_hits"`
	DirCapHits      int `json:"dir_cap_hits"`
	SymlinksSkipped int `json:"symlinks_skipped"`
}

// Summary is a derived, high-level view of a Report. Purely a function
// of Report — no I/O, no side effects.
type Summary struct {
	Total       int            `json:"total"`
	ByEcosystem map[string]int `json:"by_ecosystem"`
	ByWorkspace map[string]int `json:"by_workspace"`
	HasMonorepo bool           `json:"has_monorepo"`
}

// Summarize projects a Report onto a Summary. O(n+m) where n is the
// component count and m is the total workspace count across components;
// bounded by the scanner's dir cap upstream.
func Summarize(r *Report) Summary {
	s := Summary{
		ByEcosystem: map[string]int{},
		ByWorkspace: map[string]int{},
	}
	if r == nil {
		return s
	}
	s.Total = len(r.Components)
	for _, c := range r.Components {
		s.ByEcosystem[c.Kind]++
		for _, w := range c.Workspaces {
			s.ByWorkspace[w.Kind]++
			s.HasMonorepo = true
		}
	}
	return s
}
