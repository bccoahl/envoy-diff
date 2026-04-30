package dotenv

// IntersectOptions controls the behaviour of Intersect.
type IntersectOptions struct {
	// OnlyChanged limits the result to keys whose values differ between left and right.
	OnlyChanged bool
}

// Intersect returns only the DiffEntry items whose keys exist in BOTH left and
// right maps. This is useful for narrowing a diff to the shared key-space and
// ignoring keys that were added or removed entirely.
//
// The returned slice preserves the order of entries as produced by Diff.
func Intersect(entries []DiffEntry, opts IntersectOptions) []DiffEntry {
	out := make([]DiffEntry, 0, len(entries))
	for _, e := range entries {
		// Skip keys that exist only on one side.
		if e.Status == StatusAdded || e.Status == StatusRemoved {
			continue
		}
		if opts.OnlyChanged && e.Status != StatusChanged {
			continue
		}
		out = append(out, e)
	}
	return out
}

// IntersectKeys returns the set of keys that are present in both maps.
func IntersectKeys(left, right map[string]string) []string {
	keys := make([]string, 0)
	for k := range left {
		if _, ok := right[k]; ok {
			keys = append(keys, k)
		}
	}
	return keys
}
