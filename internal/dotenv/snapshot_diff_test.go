package dotenv_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/dotenv"
)

func makeSnapshot(entries []dotenv.DiffEntry) *dotenv.Snapshot {
	return &dotenv.Snapshot{
		LeftLabel:  "snap-left",
		RightLabel: "snap-right",
		Entries:    entries,
		Summary:    dotenv.Summarize(entries),
	}
}

func TestCompareSnapshot_NewEntry(t *testing.T) {
	snap := makeSnapshot([]dotenv.DiffEntry{
		{Key: "EXISTING", Status: dotenv.StatusUnchanged},
	})
	live := []dotenv.DiffEntry{
		{Key: "EXISTING", Status: dotenv.StatusUnchanged},
		{Key: "BRAND_NEW", Status: dotenv.StatusAdded},
	}

	res, err := dotenv.CompareSnapshot(snap, live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.New) != 1 || res.New[0].Key != "BRAND_NEW" {
		t.Errorf("expected BRAND_NEW in New, got %v", res.New)
	}
}

func TestCompareSnapshot_GoneEntry(t *testing.T) {
	snap := makeSnapshot([]dotenv.DiffEntry{
		{Key: "WAS_HERE", Status: dotenv.StatusChanged},
		{Key: "STILL_HERE", Status: dotenv.StatusUnchanged},
	})
	live := []dotenv.DiffEntry{
		{Key: "STILL_HERE", Status: dotenv.StatusUnchanged},
	}

	res, err := dotenv.CompareSnapshot(snap, live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Gone) != 1 || res.Gone[0].Key != "WAS_HERE" {
		t.Errorf("expected WAS_HERE in Gone, got %v", res.Gone)
	}
}

func TestCompareSnapshot_ShiftedStatus(t *testing.T) {
	snap := makeSnapshot([]dotenv.DiffEntry{
		{Key: "FLIPPED", Status: dotenv.StatusChanged},
	})
	live := []dotenv.DiffEntry{
		{Key: "FLIPPED", Status: dotenv.StatusUnchanged},
	}

	res, err := dotenv.CompareSnapshot(snap, live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Shifted) != 1 {
		t.Fatalf("expected 1 shifted, got %d", len(res.Shifted))
	}
	s := res.Shifted[0]
	if s.OldStatus != dotenv.StatusChanged || s.NewStatus != dotenv.StatusUnchanged {
		t.Errorf("unexpected shift: %v -> %v", s.OldStatus, s.NewStatus)
	}
}

func TestCompareSnapshot_NilSnapshot(t *testing.T) {
	_, err := dotenv.CompareSnapshot(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestCompareSnapshot_NoChanges(t *testing.T) {
	entries := []dotenv.DiffEntry{
		{Key: "STABLE", Status: dotenv.StatusUnchanged},
	}
	snap := makeSnapshot(entries)
	res, err := dotenv.CompareSnapshot(snap, entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.New) != 0 || len(res.Gone) != 0 || len(res.Shifted) != 0 {
		t.Errorf("expected no changes, got %+v", res)
	}
}
