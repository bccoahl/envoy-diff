package dotenv

import (
	"fmt"
	"io"
	"strings"
)

// OutputFormat controls how diff results are rendered.
type OutputFormat string

const (
	FormatText  OutputFormat = "text"
	FormatJSON  OutputFormat = "json"
	FormatTable OutputFormat = "table"
)

// TextFormatter writes a human-readable diff to w.
func TextFormatter(w io.Writer, entries []DiffEntry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return
	}
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			fmt.Fprintf(w, "+ %s=%s\n", e.Key, e.ValueB)
		case StatusRemoved:
			fmt.Fprintf(w, "- %s=%s\n", e.Key, e.ValueA)
		case StatusChanged:
			fmt.Fprintf(w, "~ %s: %q -> %q\n", e.Key, e.ValueA, e.ValueB)
		case StatusUnchanged:
			fmt.Fprintf(w, "  %s=%s\n", e.Key, e.ValueA)
		}
	}
}

// TableFormatter writes a Markdown-style table diff to w.
func TableFormatter(w io.Writer, entries []DiffEntry) {
	fmt.Fprintf(w, "%-4s %-30s %-30s %-30s\n", "ST", "KEY", "VALUE_A", "VALUE_B")
	fmt.Fprintln(w, strings.Repeat("-", 98))
	for _, e := range entries {
		symbol := " "
		switch e.Status {
		case StatusAdded:
			symbol = "+"
		case StatusRemoved:
			symbol = "-"
		case StatusChanged:
			symbol = "~"
		}
		fmt.Fprintf(w, "%-4s %-30s %-30s %-30s\n", symbol, e.Key, e.ValueA, e.ValueB)
	}
}
