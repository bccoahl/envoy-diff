package dotenv

import "fmt"

// Stats holds summary counts for a set of diff entries.
type Stats struct {
	Added     int
	Removed   int
	Changed   int
	Unchanged int
	Total     int
}

// Summarize computes Stats from a slice of DiffEntry.
func Summarize(entries []DiffEntry) Stats {
	var s Stats
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			s.Added++
		case StatusRemoved:
			s.Removed++
		case StatusChanged:
			s.Changed++
		case StatusUnchanged:
			s.Unchanged++
		}
	}
	s.Total = s.Added + s.Removed + s.Changed + s.Unchanged
	return s
}

// HasDiff returns true if any entries are added, removed, or changed.
func (s Stats) HasDiff() bool {
	return s.Added > 0 || s.Removed > 0 || s.Changed > 0
}

// FormatSummary returns a human-readable one-line summary string.
func (s Stats) FormatSummary() string {
	return fmt.Sprintf(
		"total: %d  added: %d  removed: %d  changed: %d  unchanged: %d",
		s.Total, s.Added, s.Removed, s.Changed, s.Unchanged,
	)
}
