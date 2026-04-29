package dotenv

import (
	"encoding/json"
	"strings"
	"testing"
)

var sampleIssues = []ValidationIssue{
	{Key: "SECRET", Message: "matches forbidden pattern", Level: ValidationError},
	{Key: "APP_ENV", Message: "key has changed", Level: ValidationWarn},
}

func TestValidationFormatter_TextContainsSymbols(t *testing.T) {
	vf := ValidationFormatter{Format: "text"}
	out, err := vf.Render(sampleIssues)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "✖") {
		t.Errorf("expected error symbol ✖ in output: %s", out)
	}
	if !strings.Contains(out, "⚠") {
		t.Errorf("expected warn symbol ⚠ in output: %s", out)
	}
}

func TestValidationFormatter_TextContainsKeys(t *testing.T) {
	vf := ValidationFormatter{Format: "text"}
	out, err := vf.Render(sampleIssues)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected key SECRET in output")
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected key APP_ENV in output")
	}
}

func TestValidationFormatter_EmptyText(t *testing.T) {
	vf := ValidationFormatter{Format: "text"}
	out, err := vf.Render(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no validation issues") {
		t.Errorf("expected no-issues message, got: %s", out)
	}
}

func TestValidationFormatter_JSONValid(t *testing.T) {
	vf := ValidationFormatter{Format: "json"}
	out, err := vf.Render(sampleIssues)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []map[string]string
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestValidationFormatter_JSONFields(t *testing.T) {
	vf := ValidationFormatter{Format: "json"}
	out, err := vf.Render(sampleIssues)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"key"`) {
		t.Errorf("expected 'key' field in JSON")
	}
	if !strings.Contains(out, `"level"`) {
		t.Errorf("expected 'level' field in JSON")
	}
	if !strings.Contains(out, `"message"`) {
		t.Errorf("expected 'message' field in JSON")
	}
}
