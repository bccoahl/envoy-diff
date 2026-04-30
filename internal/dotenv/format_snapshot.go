package dotenv

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatSnapshotDiff returns a human-readable text summary of a SnapshotDiffResult.
func FormatSnapshotDiff(r *SnapshotDiffResult) string {
	if r == nil {
		return "(no snapshot diff result)"
	}

	var sb strings.Builder

	if len(r.New) == 0 && len(r.Gone) == 0 && len(r.Shifted) == 0 {
		sb.WriteString("No changes since snapshot.\n")
		return sb.String()
	}

	if len(r.New) > 0 {
		sb.WriteString(fmt.Sprintf("New keys (%d):\n", len(r.New)))
		for _, e := range r.New {
			sb.WriteString(fmt.Sprintf("  + %s [%s]\n", e.Key, e.Status))
		}
	}

	if len(r.Gone) > 0 {
		sb.WriteString(fmt.Sprintf("Gone keys (%d):\n", len(r.Gone)))
		for _, e := range r.Gone {
			sb.WriteString(fmt.Sprintf("  - %s [%s]\n", e.Key, e.Status))
		}
	}

	if len(r.Shifted) > 0 {
		sb.WriteString(fmt.Sprintf("Status shifts (%d):\n", len(r.Shifted)))
		for _, s := range r.Shifted {
			sb.WriteString(fmt.Sprintf("  ~ %s: %s -> %s\n", s.Key, s.OldStatus, s.NewStatus))
		}
	}

	return sb.String()
}

// FormatSnapshotDiffJSON returns a JSON-encoded representation of a SnapshotDiffResult.
func FormatSnapshotDiffJSON(r *SnapshotDiffResult) (string, error) {
	if r == nil {
		r = &SnapshotDiffResult{}
	}

	type jsonOutput struct {
		New     []DiffEntry    `json:"new"`
		Gone    []DiffEntry    `json:"gone"`
		Shifted []ShiftedEntry `json:"shifted"`
	}

	out := jsonOutput{
		New:     r.New,
		Gone:    r.Gone,
		Shifted: r.Shifted,
	}

	if out.New == nil {
		out.New = []DiffEntry{}
	}
	if out.Gone == nil {
		out.Gone = []DiffEntry{}
	}
	if out.Shifted == nil {
		out.Shifted = []ShiftedEntry{}
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", fmt.Errorf("format snapshot diff: %w", err)
	}
	return string(data), nil
}
