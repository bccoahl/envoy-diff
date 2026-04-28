package dotenv

import "fmt"

// CompareOptions controls how two Sources are compared and reported.
type CompareOptions struct {
	FilterOptions FilterOptions
	Format        string // "text", "table", "json"
}

// Formatter is the common interface for all output formatters.
type Formatter interface {
	Format(entries []DiffEntry) (string, error)
}

// Compare diffs two Sources and returns formatted output.
func Compare(a, b *Source, opts CompareOptions) (string, error) {
	entries := Diff(a.Vars, b.Vars, a.Name, b.Name)
	entries = Filter(entries, opts.FilterOptions)

	f, err := resolveFormatter(opts.Format)
	if err != nil {
		return "", err
	}
	return f.Format(entries)
}

// resolveFormatter returns the Formatter for the given format name.
func resolveFormatter(format string) (Formatter, error) {
	switch format {
	case "", "text":
		return TextFormatter{}, nil
	case "table":
		return TableFormatter{}, nil
	case "json":
		return JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown format %q: choose text, table, or json", format)
	}
}
