package dotenv

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "ADDED_KEY", Status: StatusAdded, ValueA: "", ValueB: "new"},
		{Key: "REMOVED_KEY", Status: StatusRemoved, ValueA: "old", ValueB: ""},
		{Key: "CHANGED_KEY", Status: StatusChanged, ValueA: "before", ValueB: "after"},
		{Key: "SAME_KEY", Status: StatusUnchanged, ValueA: "same", ValueB: "same"},
	}
}

func TestTextFormatter_ContainsSymbols(t *testing.T) {
	var buf bytes.Buffer
	TextFormatter(&buf, sampleEntries())
	out := buf.String()

	if !strings.Contains(out, "+ ADDED_KEY") {
		t.Errorf("expected '+' for added key, got:\n%s", out)
	}
	if !strings.Contains(out, "- REMOVED_KEY") {
		t.Errorf("expected '-' for removed key, got:\n%s", out)
	}
	if !strings.Contains(out, "~ CHANGED_KEY") {
		t.Errorf("expected '~' for changed key, got:\n%s", out)
	}
}

func TestTextFormatter_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	TextFormatter(&buf, []DiffEntry{})
	if !strings.Contains(buf.String(), "No differences") {
		t.Error("expected 'No differences' message for empty entries")
	}
}

func TestTableFormatter_HasHeader(t *testing.T) {
	var buf bytes.Buffer
	TableFormatter(&buf, sampleEntries())
	out := buf.String()
	if !strings.Contains(out, "KEY") || !strings.Contains(out, "VALUE_A") {
		t.Errorf("expected table header, got:\n%s", out)
	}
}

func TestJSONFormatter_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := JSONFormatter(&buf, sampleEntries()); err != nil {
		t.Fatalf("JSONFormatter error: %v", err)
	}
	var records []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(records) != len(sampleEntries()) {
		t.Errorf("expected %d records, got %d", len(sampleEntries()), len(records))
	}
}

func TestJSONFormatter_StatusField(t *testing.T) {
	var buf bytes.Buffer
	_ = JSONFormatter(&buf, sampleEntries())
	if !strings.Contains(buf.String(), `"status"`) {
		t.Error("expected 'status' field in JSON output")
	}
}
