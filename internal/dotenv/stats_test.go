package dotenv

import (
	"strings"
	"testing"
)

var statEntries = []DiffEntry{
	{Key: "A", Status: StatusAdded, RightVal: "1"},
	{Key: "B", Status: StatusAdded, RightVal: "2"},
	{Key: "C", Status: StatusRemoved, LeftVal: "old"},
	{Key: "D", Status: StatusChanged, LeftVal: "x", RightVal: "y"},
	{Key: "E", Status: StatusUnchanged, LeftVal: "z", RightVal: "z"},
	{Key: "F", Status: StatusUnchanged, LeftVal: "w", RightVal: "w"},
}

func TestSummarize_Counts(t *testing.T) {
	s := Summarize(statEntries)
	if s.Added != 2 {
		t.Errorf("expected Added=2, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", s.Removed)
	}
	if s.Changed != 1 {
		t.Errorf("expected Changed=1, got %d", s.Changed)
	}
	if s.Unchanged != 2 {
		t.Errorf("expected Unchanged=2, got %d", s.Unchanged)
	}
	if s.Total != 6 {
		t.Errorf("expected Total=6, got %d", s.Total)
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize([]DiffEntry{})
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
	if s.HasDiff() {
		t.Error("expected HasDiff=false for empty entries")
	}
}

func TestHasDiff_True(t *testing.T) {
	s := Summarize(statEntries)
	if !s.HasDiff() {
		t.Error("expected HasDiff=true")
	}
}

func TestHasDiff_OnlyUnchanged(t *testing.T) {
	entries := []DiffEntry{
		{Key: "X", Status: StatusUnchanged, LeftVal: "v", RightVal: "v"},
	}
	s := Summarize(entries)
	if s.HasDiff() {
		t.Error("expected HasDiff=false when only unchanged entries")
	}
}

func TestFormatSummary_ContainsFields(t *testing.T) {
	s := Summarize(statEntries)
	line := s.FormatSummary()
	for _, want := range []string{"total:", "added:", "removed:", "changed:", "unchanged:"} {
		if !strings.Contains(line, want) {
			t.Errorf("FormatSummary missing %q in: %s", want, line)
		}
	}
}
