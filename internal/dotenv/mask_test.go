package dotenv

import (
	"testing"
)

func baseMaskEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "DB_PASSWORD", OldValue: "secret", NewValue: "newsecret", Status: StatusChanged},
		{Key: "API_KEY", OldValue: "", NewValue: "abc123", Status: StatusAdded},
		{Key: "APP_NAME", OldValue: "myapp", NewValue: "myapp", Status: StatusUnchanged},
		{Key: "SECRET_TOKEN", OldValue: "tok", NewValue: "", Status: StatusRemoved},
	}
}

func TestMask_MaskAll(t *testing.T) {
	entries := baseMaskEntries()
	result := Mask(entries, MaskOptions{MaskAll: true})
	for _, e := range result {
		if e.OldValue != "" && e.OldValue != "***" {
			t.Errorf("key %s: OldValue not masked, got %q", e.Key, e.OldValue)
		}
		if e.NewValue != "" && e.NewValue != "***" {
			t.Errorf("key %s: NewValue not masked, got %q", e.Key, e.NewValue)
		}
	}
}

func TestMask_ByKey(t *testing.T) {
	entries := baseMaskEntries()
	result := Mask(entries, MaskOptions{Keys: []string{"DB_PASSWORD"}})

	for _, e := range result {
		if e.Key == "DB_PASSWORD" {
			if e.OldValue != "***" || e.NewValue != "***" {
				t.Errorf("DB_PASSWORD should be masked, got old=%q new=%q", e.OldValue, e.NewValue)
			}
		} else if e.Key == "APP_NAME" {
			if e.OldValue != "myapp" {
				t.Errorf("APP_NAME should not be masked, got %q", e.OldValue)
			}
		}
	}
}

func TestMask_ByPattern(t *testing.T) {
	entries := baseMaskEntries()
	result := Mask(entries, MaskOptions{Patterns: []string{"(?i)secret|(?i)token"}})

	maskedKeys := map[string]bool{"DB_PASSWORD": false, "SECRET_TOKEN": true, "API_KEY": false, "APP_NAME": false}
	for _, e := range result {
		expectMasked := maskedKeys[e.Key]
		if expectMasked {
			if e.OldValue != "" && e.OldValue != "***" {
				t.Errorf("key %s OldValue should be masked", e.Key)
			}
		}
	}
}

func TestMask_DoesNotMutateOriginal(t *testing.T) {
	entries := baseMaskEntries()
	original := entries[0].OldValue
	Mask(entries, MaskOptions{MaskAll: true})
	if entries[0].OldValue != original {
		t.Error("Mask should not mutate the original slice")
	}
}

func TestMask_EmptyOptions(t *testing.T) {
	entries := baseMaskEntries()
	result := Mask(entries, MaskOptions{})
	for i, e := range result {
		if e.OldValue != entries[i].OldValue || e.NewValue != entries[i].NewValue {
			t.Errorf("key %s: values should be unchanged with empty options", e.Key)
		}
	}
}
