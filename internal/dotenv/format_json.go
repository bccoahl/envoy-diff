package dotenv

import (
	"encoding/json"
	"io"
)

// jsonEntry is the serialisable form of a DiffEntry.
type jsonEntry struct {
	Key    string `json:"key"`
	Status string `json:"status"`
	ValueA string `json:"value_a,omitempty"`
	ValueB string `json:"value_b,omitempty"`
}

// JSONFormatter writes diff results as a JSON array to w.
func JSONFormatter(w io.Writer, entries []DiffEntry) error {
	records := make([]jsonEntry, 0, len(entries))
	for _, e := range entries {
		records = append(records, jsonEntry{
			Key:    e.Key,
			Status: string(e.Status),
			ValueA: e.ValueA,
			ValueB: e.ValueB,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}
