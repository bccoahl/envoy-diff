package dotenv

import "sort"

// SortOrder defines the ordering strategy for diff entries.
type SortOrder string

const (
	// SortByKey sorts entries alphabetically by key name.
	SortByKey SortOrder = "key"
	// SortByStatus sorts entries grouped by their diff status (added, removed, changed, unchanged).
	SortByStatus SortOrder = "status"
	// SortNone preserves the original union-key order.
	SortNone SortOrder = "none"
)

// statusRank assigns a numeric rank to each DiffStatus for stable ordering.
func statusRank(s DiffStatus) int {
	switch s {
	case StatusAdded:
		return 0
	case StatusRemoved:
		return 1
	case StatusChanged:
		return 2
	case StatusUnchanged:
		return 3
	default:
		return 4
	}
}

// SortEntries returns a new slice of DiffEntry values sorted according to the
// given SortOrder. The original slice is not modified.
func SortEntries(entries []DiffEntry, order SortOrder) []DiffEntry {
	if order == SortNone || len(entries) == 0 {
		return entries
	}

	copied := make([]DiffEntry, len(entries))
	copy(copied, entries)

	switch order {
	case SortByKey:
		sort.SliceStable(copied, func(i, j int) bool {
			return copied[i].Key < copied[j].Key
		})
	case SortByStatus:
		sort.SliceStable(copied, func(i, j int) bool {
			ri := statusRank(copied[i].Status)
			rj := statusRank(copied[j].Status)
			if ri != rj {
				return ri < rj
			}
			return copied[i].Key < copied[j].Key
		})
	}

	return copied
}
