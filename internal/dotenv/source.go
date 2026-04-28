package dotenv

import (
	"fmt"
	"os"
	"strings"
)

// Source represents a named environment variable source.
type Source struct {
	Name string
	Vars map[string]string
}

// SourceType identifies the kind of source to load.
type SourceType string

const (
	SourceTypeFile SourceType = "file"
	SourceTypeEnv  SourceType = "env"
)

// SourceSpec describes how to load a source.
type SourceSpec struct {
	Type SourceType
	Ref  string // file path or prefix for env
}

// LoadSource loads a Source from a SourceSpec.
func LoadSource(spec SourceSpec) (*Source, error) {
	switch spec.Type {
	case SourceTypeFile:
		vars, err := LoadFile(spec.Ref)
		if err != nil {
			return nil, fmt.Errorf("load file source %q: %w", spec.Ref, err)
		}
		return &Source{Name: spec.Ref, Vars: vars}, nil
	case SourceTypeEnv:
		vars := loadFromEnvironment(spec.Ref)
		name := spec.Ref
		if name == "" {
			name = "<env>"
		}
		return &Source{Name: name, Vars: vars}, nil
	default:
		return nil, fmt.Errorf("unknown source type: %q", spec.Type)
	}
}

// loadFromEnvironment reads all OS environment variables, optionally
// filtering to those with the given prefix (prefix is stripped from keys).
func loadFromEnvironment(prefix string) map[string]string {
	vars := make(map[string]string)
	for _, entry := range os.Environ() {
		key, val, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		if prefix != "" {
			if !strings.HasPrefix(key, prefix) {
				continue
			}
			key = strings.TrimPrefix(key, prefix)
		}
		vars[key] = val
	}
	return vars
}
