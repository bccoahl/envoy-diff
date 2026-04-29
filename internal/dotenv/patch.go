package dotenv

import (
	"fmt"
	"strings"
)

// PatchFormat defines the output format for patch generation.
type PatchFormat string

const (
	PatchFormatDotenv PatchFormat = "dotenv"
	PatchFormatShell  PatchFormat = "shell"
)

// PatchOptions controls how the patch is generated.
type PatchOptions struct {
	// Format determines the output syntax.
	Format PatchFormat
	// OnlyChanged limits the patch to added/changed/removed entries.
	OnlyChanged bool
	// Side is either "left" or "right" (default "right").
	Side string
}

// Patch generates a patch script that, when applied to the left environment,
// produces the right environment. It returns lines ready to write to a file.
func Patch(entries []DiffEntry, opts PatchOptions) ([]string, error) {
	if opts.Format == "" {
		opts.Format = PatchFormatDotenv
	}
	if opts.Side == "" {
		opts.Side = "right"
	}
	if opts.Side != "left" && opts.Side != "right" {
		return nil, fmt.Errorf("patch: unknown side %q, must be \"left\" or \"right\"", opts.Side)
	}

	var lines []string

	for _, e := range entries {
		if opts.OnlyChanged && e.Status == StatusUnchanged {
			continue
		}

		var value string
		var remove bool

		switch {
		case e.Status == StatusRemoved && opts.Side == "right":
			remove = true
		case e.Status == StatusAdded && opts.Side == "left":
			remove = true
		default:
			if opts.Side == "right" {
				value = e.RightValue
			} else {
				value = e.LeftValue
			}
		}

		switch opts.Format {
		case PatchFormatShell:
			if remove {
				lines = append(lines, fmt.Sprintf("unset %s", e.Key))
			} else {
				lines = append(lines, fmt.Sprintf("export %s=%s", e.Key, shellQuoteValue(value)))
			}
		default: // dotenv
			if remove {
				lines = append(lines, fmt.Sprintf("# REMOVE %s", e.Key))
			} else {
				lines = append(lines, fmt.Sprintf("%s=%s", e.Key, dotenvQuoteValue(value)))
			}
		}
	}

	return lines, nil
}

func shellQuoteValue(v string) string {
	if !strings.ContainsAny(v, " \t\n\"'\\$`") {
		return v
	}
	return "'" + strings.ReplaceAll(v, "'", "'\"'\"'") + "'"
}

func dotenvQuoteValue(v string) string {
	if !strings.ContainsAny(v, " \t\n\"'") {
		return v
	}
	return `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
}
