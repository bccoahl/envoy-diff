package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of a diff result with metadata.
type Snapshot struct {
	CreatedAt time.Time   `json:"created_at"`
	LeftLabel  string      `json:"left_label"`
	RightLabel string      `json:"right_label"`
	Entries    []DiffEntry `json:"entries"`
	Summary    Summary     `json:"summary"`
}

// SaveSnapshot writes a Snapshot to the given file path as JSON.
func SaveSnapshot(path, leftLabel, rightLabel string, entries []DiffEntry) error {
	snap := Snapshot{
		CreatedAt:  time.Now().UTC(),
		LeftLabel:  leftLabel,
		RightLabel: rightLabel,
		Entries:    entries,
		Summary:    Summarize(entries),
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// LoadSnapshot reads a Snapshot from the given JSON file path.
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: parse failed: %w", err)
	}

	return &snap, nil
}
