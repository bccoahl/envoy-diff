package dotenv

import (
	"strings"
	"testing"
)

var baseValidateEntries = []DiffEntry{
	{Key: "APP_ENV", LeftVal: "production", RightVal: "staging", Status: StatusChanged},
	{Key: "DB_HOST", LeftVal: "db.prod", RightVal: "db.prod", Status: StatusUnchanged},
	{Key: "SECRET_KEY", LeftVal: "", RightVal: "abc123", Status: StatusAdded},
	{Key: "OLD_FLAG", LeftVal: "true", RightVal: "", Status: StatusRemoved},
}

func TestValidate_RequireKeys_Present(t *testing.T) {
	issues := Validate(baseValidateEntries, ValidateOptions{
		RequireKeys: []string{"APP_ENV", "DB_HOST"},
	})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestValidate_RequireKeys_Missing(t *testing.T) {
	issues := Validate(baseValidateEntries, ValidateOptions{
		RequireKeys: []string{"MISSING_KEY"},
	})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Level != ValidationError {
		t.Errorf("expected error level, got %s", issues[0].Level)
	}
	if issues[0].Key != "MISSING_KEY" {
		t.Errorf("unexpected key: %s", issues[0].Key)
	}
}

func TestValidate_ForbidPattern(t *testing.T) {
	issues := Validate(baseValidateEntries, ValidateOptions{
		ForbidPattern: "^SECRET_",
	})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY, got %s", issues[0].Key)
	}
	if issues[0].Level != ValidationError {
		t.Errorf("expected error level")
	}
}

func TestValidate_WarnOnChanged(t *testing.T) {
	issues := Validate(baseValidateEntries, ValidateOptions{
		WarnOnChanged: true,
	})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Level != ValidationWarn {
		t.Errorf("expected warn level, got %s", issues[0].Level)
	}
	if issues[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV, got %s", issues[0].Key)
	}
}

func TestHasErrors_True(t *testing.T) {
	issues := []ValidationIssue{
		{Key: "X", Message: "bad", Level: ValidationError},
	}
	if !HasErrors(issues) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_WarnOnly(t *testing.T) {
	issues := []ValidationIssue{
		{Key: "X", Message: "warn", Level: ValidationWarn},
	}
	if HasErrors(issues) {
		t.Error("expected HasErrors to return false for warn-only issues")
	}
}

func TestFormatIssues_Empty(t *testing.T) {
	out := FormatIssues(nil)
	if !strings.Contains(out, "no validation issues") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatIssues_ContainsLevel(t *testing.T) {
	issues := []ValidationIssue{
		{Key: "FOO", Message: "something wrong", Level: ValidationError},
		{Key: "BAR", Message: "heads up", Level: ValidationWarn},
	}
	out := FormatIssues(issues)
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected [ERROR] in output: %s", out)
	}
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected [WARN] in output: %s", out)
	}
}
