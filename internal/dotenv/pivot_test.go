package dotenv

import (
	"testing"
)

func sources(pairs ...interface{}) []NamedSource {
	var ns []NamedSource
	for i := 0; i+1 < len(pairs); i += 2 {
		ns = append(ns, NamedSource{
			Name: pairs[i].(string),
			Env:  pairs[i+1].(map[string]string),
		})
	}
	return ns
}

func findPivot(entries []PivotEntry, key string) (PivotEntry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return PivotEntry{}, false
}

func TestPivot_AllPresent_SameValue(t *testing.T) {
	srcs := sources(
		"prod", map[string]string{"HOST": "example.com", "PORT": "443"},
		"staging", map[string]string{"HOST": "example.com", "PORT": "443"},
	)
	entries := Pivot(srcs, PivotOptions{})
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	e, ok := findPivot(entries, "HOST")
	if !ok {
		t.Fatal("HOST not found")
	}
	if e.Values["prod"] != "example.com" || e.Values["staging"] != "example.com" {
		t.Errorf("unexpected values: %v", e.Values)
	}
	if len(e.Missing) != 0 {
		t.Errorf("expected no missing, got %v", e.Missing)
	}
}

func TestPivot_KeyMissingInOneSource(t *testing.T) {
	srcs := sources(
		"prod", map[string]string{"HOST": "prod.example.com", "SECRET": "s3cr3t"},
		"staging", map[string]string{"HOST": "staging.example.com"},
	)
	entries := Pivot(srcs, PivotOptions{})
	e, ok := findPivot(entries, "SECRET")
	if !ok {
		t.Fatal("SECRET not found")
	}
	if len(e.Missing) != 1 || e.Missing[0] != "staging" {
		t.Errorf("expected staging in Missing, got %v", e.Missing)
	}
}

func TestPivot_OnlyDiff_ExcludesIdentical(t *testing.T) {
	srcs := sources(
		"a", map[string]string{"SAME": "x", "DIFF": "1"},
		"b", map[string]string{"SAME": "x", "DIFF": "2"},
	)
	entries := Pivot(srcs, PivotOptions{OnlyDiff: true})
	if _, ok := findPivot(entries, "SAME"); ok {
		t.Error("SAME should be excluded when OnlyDiff=true")
	}
	if _, ok := findPivot(entries, "DIFF"); !ok {
		t.Error("DIFF should be included when OnlyDiff=true")
	}
}

func TestPivot_KeysAreSorted(t *testing.T) {
	srcs := sources(
		"env", map[string]string{"ZEBRA": "z", "ALPHA": "a", "MIDDLE": "m"},
	)
	entries := Pivot(srcs, PivotOptions{})
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, e := range entries {
		if e.Key != expected[i] {
			t.Errorf("position %d: want %s, got %s", i, expected[i], e.Key)
		}
	}
}

func TestPivot_EmptySources(t *testing.T) {
	entries := Pivot([]NamedSource{}, PivotOptions{})
	if len(entries) != 0 {
		t.Errorf("expected empty result, got %d entries", len(entries))
	}
}
