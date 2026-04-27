package dotenv

import "strings"

// FilterOptions controls which diff entries are included in output.
type FilterOptions struct {
	// OnlyChanged includes only Added, Removed, and Changed entries.
	OnlyChanged bool
	// Keys restricts output to entries whose key matches one of the provided keys.
	// If empty, all keys are included.
	Keys []string
	// Prefix restricts output to entries whose key starts with the given prefix.
	Prefix string
}

// Filter returns a subset of entries based on the provided FilterOptions.
func Filter(entries []DiffEntry, opts FilterOptions) []DiffEntry {
	var result []DiffEntry

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	for _, e := range entries {
		if opts.OnlyChanged && e.Status == StatusUnchanged {
			continue
		}

		if len(opts.Keys) > 0 {
			if _, ok := keySet[e.Key]; !ok {
				continue
			}
		}

		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}

		result = append(result, e)
	}

	return result
}
