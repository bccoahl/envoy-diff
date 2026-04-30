package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envoy-diff/internal/dotenv"
)

var baseSnapshotEntries = []dotenv.DiffEntry{
	{Key: "APP_ENV", Left: "staging", Right: "production", Status: dotenv.StatusChanged},
	{Key: "DB_HOST", Left: "localhost", Right: "localhost", Status: dotenv.StatusUnchanged},
	{Key: "NEW_KEY", Left: "", Right: "value", Status: dotenv.StatusAdded},
	{Key: "OLD_KEY", Left: "old", Right: "", Status: dotenv.StatusRemoved},
}

func TestSaveSnapshot_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	err := dotenv.SaveSnapshot(path, "left.env", "right.env", baseSnapshotEntries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected snapshot file to exist")
	}
}

func TestLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	before := time.Now().UTC().Truncate(time.Second)
	err := dotenv.SaveSnapshot(path, "a.env", "b.env", baseSnapshotEntries)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	snap, err := dotenv.LoadSnapshot(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	if snap.LeftLabel != "a.env" {
		t.Errorf("expected left label 'a.env', got %q", snap.LeftLabel)
	}
	if snap.RightLabel != "b.env" {
		t.Errorf("expected right label 'b.env', got %q", snap.RightLabel)
	}
	if len(snap.Entries) != len(baseSnapshotEntries) {
		t.Errorf("expected %d entries, got %d", len(baseSnapshotEntries), len(snap.Entries))
	}
	if snap.CreatedAt.Before(before) {
		t.Error("expected CreatedAt to be after test start")
	}
}

func TestLoadSnapshot_SummaryPopulated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	_ = dotenv.SaveSnapshot(path, "x", "y", baseSnapshotEntries)
	snap, _ := dotenv.LoadSnapshot(path)

	if snap.Summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", snap.Summary.Added)
	}
	if snap.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", snap.Summary.Removed)
	}
	if snap.Summary.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", snap.Summary.Changed)
	}
}

func TestLoadSnapshot_FileNotFound(t *testing.T) {
	_, err := dotenv.LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)

	_, err := dotenv.LoadSnapshot(path)
	if err == nil {
		t.Fatal("expected parse error for invalid JSON")
	}
}
