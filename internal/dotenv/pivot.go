package dotenv

// PivotEntry represents a single key compared across multiple named sources.
type PivotEntry struct {
	Key    string
	Values map[string]string // source name -> value
	Missing []string          // source names where key is absent
}

// PivotOptions controls the behaviour of Pivot.
type PivotOptions struct {
	// OnlyDiff, when true, excludes keys whose values are identical across all sources.
	OnlyDiff bool
}

// Pivot compares multiple named environment maps and returns a per-key view
// that shows each source's value side-by-side.
//
// sources is an ordered slice of (name, map) pairs so that column order is
// deterministic.
func Pivot(sources []NamedSource, opts PivotOptions) []PivotEntry {
	// Collect union of all keys.
	keySeen := make(map[string]struct{})
	for _, ns := range sources {
		for k := range ns.Env {
			keySeen[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySeen))
	for k := range keySeen {
		keys = append(keys, k)
	}
	sortStrings(keys)

	var entries []PivotEntry
	for _, k := range keys {
		entry := PivotEntry{
			Key:    k,
			Values: make(map[string]string, len(sources)),
		}
		for _, ns := range sources {
			if v, ok := ns.Env[k]; ok {
				entry.Values[ns.Name] = v
			} else {
				entry.Missing = append(entry.Missing, ns.Name)
			}
		}
		if opts.OnlyDiff && !pivotHasDiff(entry, sources) {
			continue
		}
		entries = append(entries, entry)
	}
	return entries
}

// NamedSource pairs a human-readable name with a flat env map.
type NamedSource struct {
	Name string
	Env  map[string]string
}

// pivotHasDiff returns true if the key has different values across sources
// or is missing from at least one source.
func pivotHasDiff(e PivotEntry, sources []NamedSource) bool {
	if len(e.Missing) > 0 {
		return true
	}
	var ref string
	first := true
	for _, ns := range sources {
		v := e.Values[ns.Name]
		if first {
			ref = v
			first = false
			continue
		}
		if v != ref {
			return true
		}
	}
	return false
}

// sortStrings sorts a string slice in-place (avoids importing sort at package
// level — reuses the helper already available via sort.go).
func sortStrings(ss []string) {
	// insertion sort — input sets are small in practice.
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
