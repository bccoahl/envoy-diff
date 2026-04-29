package dotenv

import (
	"fmt"
	"os"
	"strings"
)

// ResolveOptions controls how variable references are resolved.
type ResolveOptions struct {
	// Expand resolves ${VAR} and $VAR references within values.
	Expand bool
	// Strict returns an error if a referenced variable is not found.
	Strict bool
}

// Resolve expands variable references in the given env map.
// It processes each value using os.Expand with the map as the lookup source.
// Earlier keys are available to later ones within the same pass.
func Resolve(env map[string]string, opts ResolveOptions) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = v
	}

	if !opts.Expand {
		return result, nil
	}

	var expandErr error
	for k, v := range result {
		if !strings.ContainsAny(v, "$") {
			continue
		}
		expanded := os.Expand(v, func(ref string) string {
			if val, ok := result[ref]; ok {
				return val
			}
			if opts.Strict {
				expandErr = fmt.Errorf("resolve: undefined variable %q referenced in key %q", ref, k)
			}
			return ""
		})
		if expandErr != nil {
			return nil, expandErr
		}
		result[k] = expanded
	}
	return result, nil
}
