package dotenv

import "strings"

// ResolveOptions controls how variable references are resolved.
type ResolveOptions struct {
	// Expand resolves ${VAR} and $VAR references within values.
	Expand bool
	// Strict returns an error if a referenced variable is not found.
	Strict bool
}

// Resolve expands variable references in the given env map.
// It processes entries in the order they are provided, so earlier
// keys can be referenced by later ones.
func Resolve(env map[string]string, opts ResolveOptions) (map[string]string, error) {
	if !opts.Expand {
		result := make(map[string]string, len(env))
		for k, v := range env {
			result[k] = v
		}
		return result, nil
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = v
	}

	var expandErr error
	for k, v := range result {
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

	_ = strings.TrimSpace // keep import
	return result, nil
}
