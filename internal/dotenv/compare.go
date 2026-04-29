package dotenv

import (
	"fmt"
	"io"
)

// CompareOptions holds configuration for a full diff + format pipeline run.
type CompareOptions struct {
	FilterOptions FilterOptions
	MaskOptions   MaskOptions
	Format        string
	SortOrder     SortOrder
}

// resolveFormatter returns the Formatter matching the requested format string.
func resolveFormatter(format string) (Formatter, error) {
	switch format {
	case "text", "":
		return &TextFormatter{}, nil
	case "table":
		return &TableFormatter{}, nil
	case "json":
		return &JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown format %q: must be one of text, table, json", format)
	}
}

// Compare diffs two env maps, applies filtering, masking, sorting, and writes
// the formatted result to w.
func Compare(a, b map[string]string, opts CompareOptions, w io.Writer) error {
	entries := Diff(a, b)

	entries = Filter(entries, opts.FilterOptions)
	entries = Mask(entries, opts.MaskOptions)
	entries = SortEntries(entries, opts.SortOrder)

	fmt, err := resolveFormatter(opts.Format)
	if err != nil {
		return err
	}

	return fmt.Format(w, entries)
}
