package dotenv_test

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/dotenv"
)

func baseIntersectEntries() []dotenv.DiffEntry {
	return []dotenv.DiffEntry{
		{Key: "SHARED_SAME", LeftValue: "foo", RightValue: "foo", Status: dotenv.StatusUnchanged},
		{Key: "SHARED_DIFF", LeftValue: "old", RightValue: "new", Status: dotenv.StatusChanged},
		{Key: "ONLY_LEFT", LeftValue: "bar", RightValue: "", Status: dotenv.StatusRemoved},
		{Key: "ONLY_RIGHT", LeftValue: "", RightValue: "baz", Status: dotenv.StatusAdded},
	}
}

func TestIntersect_ExcludesAddedAndRemoved(t *testing.T) {
	result := dotenv.Intersect(baseIntersectEntries(), dotenv.IntersectOptions{})
	for _, e := range result {
		if e.Status == dotenv.StatusAdded || e.Status == dotenv.StatusRemoved {
			t.Errorf("unexpected status %s for key %s", e.Status, e.Key)
		}
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestIntersect_OnlyChanged(t *testing.T) {
	result := dotenv.Intersect(baseIntersectEntries(), dotenv.IntersectOptions{OnlyChanged: true})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Key != "SHARED_DIFF" {
		t.Errorf("expected SHARED_DIFF, got %s", result[0].Key)
	}
}

func TestIntersect_Empty(t *testing.T) {
	result := dotenv.Intersect([]dotenv.DiffEntry{}, dotenv.IntersectOptions{})
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}

func TestIntersect_AllAddedRemoved(t *testing.T) {
	entries := []dotenv.DiffEntry{
		{Key: "A", Status: dotenv.StatusAdded},
		{Key: "B", Status: dotenv.StatusRemoved},
	}
	result := dotenv.Intersect(entries, dotenv.IntersectOptions{})
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}

func TestIntersectKeys_ReturnsCommonKeys(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2", "C": "3"}
	right := map[string]string{"B": "2", "C": "99", "D": "4"}
	keys := dotenv.IntersectKeys(left, right)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(keys), keys)
	}
	keySet := map[string]bool{}
	for _, k := range keys {
		keySet[k] = true
	}
	if !keySet["B"] || !keySet["C"] {
		t.Errorf("expected B and C in result, got %v", keys)
	}
}

func TestIntersectKeys_NoOverlap(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"B": "2"}
	keys := dotenv.IntersectKeys(left, right)
	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}
