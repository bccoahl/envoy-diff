package dotenv

import (
	"testing"
)

var baseRedactEntries = []DiffEntry{
	{Key: "DATABASE_URL", OldValue: "postgres://old", NewValue: "postgres://new", Status: StatusChanged},
	{Key: "API_KEY", OldValue: "secret-old", NewValue: "secret-new", Status: StatusChanged},
	{Key: "APP_ENV", OldValue: "staging", NewValue: "production", Status: StatusChanged},
	{Key: "AWS_SECRET_ACCESS_KEY", OldValue: "abc123", NewValue: "abc123", Status: StatusUnchanged},
	{Key: "PORT", OldValue: "8080", NewValue: "9090", Status: StatusChanged},
}

func TestRedact_ByKey(t *testing.T) {
	result := Redact(baseRedactEntries, RedactOptions{
		Keys: []string{"API_KEY"},
	})
	for _, e := range result {
		if e.Key == "API_KEY" {
			if e.OldValue != "[REDACTED]" || e.NewValue != "[REDACTED]" {
				t.Errorf("expected API_KEY to be redacted, got old=%q new=%q", e.OldValue, e.NewValue)
			}
		} else if e.Key == "DATABASE_URL" {
			if e.OldValue == "[REDACTED]" {
				t.Errorf("DATABASE_URL should not be redacted")
			}
		}
	}
}

func TestRedact_ByPattern(t *testing.T) {
	result := Redact(baseRedactEntries, RedactOptions{
		Patterns: []string{`(?i)secret|password|key`},
	})
	redacted := map[string]bool{}
	for _, e := range result {
		if e.OldValue == "[REDACTED]" {
			redacted[e.Key] = true
		}
	}
	if !redacted["API_KEY"] {
		t.Error("expected API_KEY to be redacted by pattern")
	}
	if !redacted["AWS_SECRET_ACCESS_KEY"] {
		t.Error("expected AWS_SECRET_ACCESS_KEY to be redacted by pattern")
	}
	if redacted["PORT"] {
		t.Error("PORT should not be redacted")
	}
}

func TestRedact_CustomReplacement(t *testing.T) {
	result := Redact(baseRedactEntries, RedactOptions{
		Keys:        []string{"DATABASE_URL"},
		Replacement: "***",
	})
	for _, e := range result {
		if e.Key == "DATABASE_URL" {
			if e.OldValue != "***" || e.NewValue != "***" {
				t.Errorf("expected custom replacement, got %q / %q", e.OldValue, e.NewValue)
			}
		}
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	orig := make([]DiffEntry, len(baseRedactEntries))
	copy(orig, baseRedactEntries)

	Redact(baseRedactEntries, RedactOptions{
		Keys: []string{"API_KEY", "DATABASE_URL"},
	})

	for i, e := range baseRedactEntries {
		if e.OldValue != orig[i].OldValue || e.NewValue != orig[i].NewValue {
			t.Errorf("original entry %q was mutated", e.Key)
		}
	}
}

func TestRedact_EmptyOptions(t *testing.T) {
	result := Redact(baseRedactEntries, RedactOptions{})
	for i, e := range result {
		if e.OldValue != baseRedactEntries[i].OldValue {
			t.Errorf("expected no redaction with empty options, key=%q", e.Key)
		}
	}
}
