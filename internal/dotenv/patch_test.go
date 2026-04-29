package dotenv

import (
	"strings"
	"testing"
)

func basePatchEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "HOST", LeftValue: "localhost", RightValue: "prod.example.com", Status: StatusChanged},
		{Key: "PORT", LeftValue: "8080", RightValue: "8080", Status: StatusUnchanged},
		{Key: "NEW_KEY", LeftValue: "", RightValue: "new_value", Status: StatusAdded},
		{Key: "OLD_KEY", LeftValue: "old_value", RightValue: "", Status: StatusRemoved},
	}
}

func TestPatch_DotenvDefault(t *testing.T) {
	lines, err := Patch(basePatchEntries(), PatchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "HOST=prod.example.com") {
		t.Errorf("expected changed key in output, got:\n%s", joined)
	}
	if !strings.Contains(joined, "NEW_KEY=new_value") {
		t.Errorf("expected added key in output, got:\n%s", joined)
	}
	if !strings.Contains(joined, "# REMOVE OLD_KEY") {
		t.Errorf("expected remove comment for OLD_KEY, got:\n%s", joined)
	}
}

func TestPatch_OnlyChanged(t *testing.T) {
	lines, err := Patch(basePatchEntries(), PatchOptions{OnlyChanged: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, l := range lines {
		if strings.HasPrefix(l, "PORT=") {
			t.Errorf("unchanged key PORT should not appear in patch")
		}
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (changed+added+removed), got %d", len(lines))
	}
}

func TestPatch_ShellFormat(t *testing.T) {
	lines, err := Patch(basePatchEntries(), PatchOptions{Format: PatchFormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "export HOST=prod.example.com") {
		t.Errorf("expected export statement, got:\n%s", joined)
	}
	if !strings.Contains(joined, "unset OLD_KEY") {
		t.Errorf("expected unset for removed key, got:\n%s", joined)
	}
}

func TestPatch_LeftSide(t *testing.T) {
	lines, err := Patch(basePatchEntries(), PatchOptions{Side: "left"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "HOST=localhost") {
		t.Errorf("expected left value for HOST, got:\n%s", joined)
	}
	if !strings.Contains(joined, "# REMOVE NEW_KEY") {
		t.Errorf("expected remove for added key on left side, got:\n%s", joined)
	}
}

func TestPatch_InvalidSide(t *testing.T) {
	_, err := Patch(basePatchEntries(), PatchOptions{Side: "both"})
	if err == nil {
		t.Error("expected error for invalid side, got nil")
	}
}

func TestPatch_QuotedValues(t *testing.T) {
	entries := []DiffEntry{
		{Key: "MSG", LeftValue: "", RightValue: "hello world", Status: StatusAdded},
	}
	lines, err := Patch(entries, PatchOptions{Format: PatchFormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 1 || !strings.Contains(lines[0], `"hello world"`) {
		t.Errorf("expected quoted value, got: %v", lines)
	}
}
