package dotenv

import "sort"

// DiffStatus represents the kind of change for a key.
type DiffStatus string

const (
	StatusAdded     DiffStatus = "added"
	StatusRemoved   DiffStatus = "removed"
	StatusChanged   DiffStatus = "changed"
	StatusUnchanged DiffStatus = "unchanged"
)

// DiffEntry holds the comparison result for a single key.
type DiffEntry struct {
	Key    string
	Status DiffStatus
	ValueA string // value in the first (left) set
	ValueB string // value in the second (right) set
}

// Diff compares two env maps and returns a sorted slice of DiffEntry.
func Diff(a, b map[string]string) []DiffEntry {
	keys := unionKeys(a, b)
	sort.Strings(keys)

	entries := make([]DiffEntry, 0, len(keys))
	for _, k := range keys {
		va, inA := a[k]
		vb, inB := b[k]

		var status DiffStatus
		switch {
		case inA && !inB:
			status = StatusRemoved
		case !inA && inB:
			status = StatusAdded
		case va != vb:
			status = StatusChanged
		default:
			status = StatusUnchanged
		}

		entries = append(entries, DiffEntry{
			Key:    k,
			Status: status,
			ValueA: va,
			ValueB: vb,
		})
	}
	return entries
}

// unionKeys returns the deduplicated union of keys from both maps.
func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
