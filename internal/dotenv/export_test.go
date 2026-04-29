package dotenv

import (
	"bytes"
	"strings"
	"testing"
)

var exportEntries = []DiffEntry{
	{Key: "APP_ENV", ValA: "staging", ValB: "production", Status: StatusChanged},
	{Key: "DEBUG", ValA: "true", ValB: "true", Status: StatusUnchanged},
	{Key: "NEW_KEY", ValA: "", ValB: "hello world", Status: StatusAdded},
	{Key: "OLD_KEY", ValA: "gone", ValB: "", Status: StatusRemoved},
}

func TestExport_DotenvDefault(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, exportEntries, ExportOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production, got:\n%s", out)
	}
	if !strings.Contains(out, "DEBUG=true") {
		t.Errorf("expected DEBUG=true, got:\n%s", out)
	}
}

func TestExport_ShellFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, exportEntries, ExportOptions{Format: ExportFormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_ENV=production") {
		t.Errorf("expected shell export line, got:\n%s", out)
	}
	if !strings.Contains(out, "export NEW_KEY='hello world'") {
		t.Errorf("expected quoted value for NEW_KEY, got:\n%s", out)
	}
}

func TestExport_DockerFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, exportEntries, ExportOptions{Format: ExportFormatDocker})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "-e APP_ENV=production") {
		t.Errorf("expected docker -e flag, got:\n%s", out)
	}
}

func TestExport_OnlyChanged(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, exportEntries, ExportOptions{OnlyChanged: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "DEBUG") {
		t.Errorf("unchanged key DEBUG should be excluded, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("changed key APP_ENV should be included, got:\n%s", out)
	}
}

func TestExport_LeftSide(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, exportEntries, ExportOptions{Side: "left"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=staging") {
		t.Errorf("expected left-side value staging, got:\n%s", out)
	}
}

func TestShellQuote(t *testing.T) {
	cases := []struct{ in, want string }{
		{"", `""`},
		{"simple", "simple"},
		{"hello world", "'hello world'"},
		{"it's", `'it'\''s'`},
	}
	for _, c := range cases {
		got := shellQuote(c.in)
		if got != c.want {
			t.Errorf("shellQuote(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
