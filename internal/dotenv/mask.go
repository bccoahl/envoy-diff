package dotenv

import (
	"regexp"
	"strings"
)

// MaskOptions controls which keys should have their values masked.
type MaskOptions struct {
	// MaskAll replaces all values with "***".
	MaskAll bool
	// Keys is an explicit list of key names to mask.
	Keys []string
	// Patterns is a list of regex patterns; matching key names will be masked.
	Patterns []string
}

// Mask replaces sensitive values in DiffEntry slices according to MaskOptions.
// The original slice is not modified; a new slice is returned.
func Mask(entries []DiffEntry, opts MaskOptions) []DiffEntry {
	const placeholder = "***"

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
		if shouldMask(e.Key, opts.MaskAll, keySet, compiled) {
			if e.OldValue != "" {
				e.OldValue = placeholder
			}
			if e.NewValue != "" {
				e.NewValue = placeholder
			}
		}
		result[i] = e
	}
	return result
}

func shouldMask(key string, maskAll bool, keySet map[string]struct{}, patterns []*regexp.Regexp) bool {
	if maskAll {
		return true
	}
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
