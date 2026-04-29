package dotenv

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ExportFormat defines the output format for exported variables.
type ExportFormat string

const (
	ExportFormatShell ExportFormat = "shell"
	ExportFormatDotenv ExportFormat = "dotenv"
	ExportFormatDocker ExportFormat = "docker"
)

// ExportOptions controls how entries are exported.
type ExportOptions struct {
	Format    ExportFormat
	// OnlyChanged restricts output to added or changed entries only.
	OnlyChanged bool
	// Side selects which value to export: "left" or "right".
	Side string
}

// Export writes the selected entries to w in the requested format.
// Side selects which value to use: "left" uses KeyA/ValA, "right" (default) uses KeyB/ValB.
func Export(w io.Writer, entries []DiffEntry, opts ExportOptions) error {
	if opts.Format == "" {
		opts.Format = ExportFormatDotenv
	}
	if opts.Side == "" {
		opts.Side = "right"
	}

	filtered := make([]DiffEntry, 0, len(entries))
	for _, e := range entries {
		if opts.OnlyChanged && (e.Status == StatusUnchanged || e.Status == StatusRemoved) {
			continue
		}
		filtered = append(filtered, e)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Key < filtered[j].Key
	})

	for _, e := range filtered {
		val := e.ValB
		if opts.Side == "left" {
			val = e.ValA
		}
		var line string
		switch opts.Format {
		case ExportFormatShell:
			line = fmt.Sprintf("export %s=%s", e.Key, shellQuote(val))
		case ExportFormatDocker:
			line = fmt.Sprintf("-e %s=%s", e.Key, shellQuote(val))
		default: // dotenv
			line = fmt.Sprintf("%s=%s", e.Key, val)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

// shellQuote wraps val in single quotes if it contains special characters.
func shellQuote(val string) string {
	if val == "" {
		return `""`
	}
	specials := " \t\n$\\\"'`!{}()|&;<>"
	if strings.ContainsAny(val, specials) {
		return "'" + strings.ReplaceAll(val, "'", `'\''`) + "'"
	}
	return val
}
