package dotenv

import "sort"

// DiffKind categorises a single entry in a diff result.
type DiffKind string

const (
	Added    DiffKind = "added"    // key present in b but not in a
	Removed  DiffKind = "removed"  // key present in a but not in b
	Changed  DiffKind = "changed"  // key present in both but values differ
	Unchanged DiffKind = "unchanged" // key present in both with same value
)

// DiffEntry represents a single key comparison result.
type DiffEntry struct {
	Key      string
	Kind     DiffKind
	ValueA   string // value in the "left" / source map
	ValueB   string // value in the "right" / target map
}

// Diff compares two EnvMaps and returns a sorted slice of DiffEntry.
// All keys from both maps are represented in the output.
func Diff(a, b EnvMap) []DiffEntry {
	keys := unionKeys(a, b)
	sort.Strings(keys)

	result := make([]DiffEntry, 0, len(keys))
	for _, k := range keys {
		va, inA := a[k]
		vb, inB := b[k]

		var kind DiffKind
		switch {
		case inA && !inB:
			kind = Removed
		case !inA && inB:
			kind = Added
		case va == vb:
			kind = Unchanged
		default:
			kind = Changed
		}

		result = append(result, DiffEntry{Key: k, Kind: kind, ValueA: va, ValueB: vb})
	}
	return result
}

func unionKeys(a, b EnvMap) []string {
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
