package dotenv

import (
	"regexp"
	"strings"
)

// RedactOptions controls how values are redacted in diff entries.
type RedactOptions struct {
	// Keys is a list of exact key names to redact.
	Keys []string
	// Patterns is a list of glob-style patterns (using regexp) to match key names.
	Patterns []string
	// Replacement is the string used in place of redacted values. Defaults to "[REDACTED]".
	Replacement string
}

// Redact returns a new slice of DiffEntry with sensitive values replaced.
// It does not mutate the original entries.
func Redact(entries []DiffEntry, opts RedactOptions) []DiffEntry {
	replacement := opts.Replacement
	if replacement == "" {
		replacement = "[REDACTED]"
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = struct{}{}
	}

	var compiled []*regexp.Regexp
	for _, p := range opts.Patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}

	result := make([]DiffEntry, len(entries))
	for i, e := range entries {
		if shouldRedact(e.Key, keySet, compiled) {
			e.OldValue = replacement
			e.NewValue = replacement
		}
		result[i] = e
	}
	return result
}

func shouldRedact(key string, keySet map[string]struct{}, patterns []*regexp.Regexp) bool {
	upper := strings.ToUpper(key)
	if _, ok := keySet[upper]; ok {
		return true
	}
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
