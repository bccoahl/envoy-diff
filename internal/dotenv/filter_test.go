package dotenv

import (
	"testing"
)

var filterSample = []DiffEntry{
	{Key: "APP_NAME", OldValue: "foo", NewValue: "foo", Status: StatusUnchanged},
	{Key: "APP_ENV", OldValue: "dev", NewValue: "prod", Status: StatusChanged},
	{Key: "DB_HOST", OldValue: "", NewValue: "localhost", Status: StatusAdded},
	{Key: "SECRET_KEY", OldValue: "abc", NewValue: "", Status: StatusRemoved},
	{Key: "DB_PORT", OldValue: "5432", NewValue: "5432", Status: StatusUnchanged},
}

func TestFilter_OnlyChanged(t *testing.T) {
	result := Filter(filterSample, FilterOptions{OnlyChanged: true})
	for _, e := range result {
		if e.Status == StatusUnchanged {
			t.Errorf("expected no unchanged entries, got key %q", e.Key)
		}
	}
	if len(result) != 3 {
		t.Errorf("expected 3 changed entries, got %d", len(result))
	}
}

func TestFilter_ByKeys(t *testing.T) {
	result := Filter(filterSample, FilterOptions{Keys: []string{"APP_NAME", "DB_HOST"}})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %q", result[0].Key)
	}
	if result[1].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", result[1].Key)
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	result := Filter(filterSample, FilterOptions{Prefix: "DB_"})
	if len(result) != 2 {
		t.Fatalf("expected 2 DB_ entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key != "DB_HOST" && e.Key != "DB_PORT" {
			t.Errorf("unexpected key %q with DB_ prefix filter", e.Key)
		}
	}
}

func TestFilter_OnlyChanged_WithPrefix(t *testing.T) {
	result := Filter(filterSample, FilterOptions{OnlyChanged: true, Prefix: "APP_"})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV, got %q", result[0].Key)
	}
}

func TestFilter_EmptyOptions(t *testing.T) {
	result := Filter(filterSample, FilterOptions{})
	if len(result) != len(filterSample) {
		t.Errorf("expected all %d entries, got %d", len(filterSample), len(result))
	}
}

func TestFilter_NoMatch(t *testing.T) {
	result := Filter(filterSample, FilterOptions{Prefix: "NONEXISTENT_"})
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}
