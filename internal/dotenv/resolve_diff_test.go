package dotenv

import (
	"testing"
)

func TestResolvedDiff_ExpandsValues(t *testing.T) {
	leftEnv := map[string]string{
		"BASE": "v1",
		"URL":  "http://$BASE",
	}
	rightEnv := map[string]string{
		"BASE": "v2",
		"URL":  "http://$BASE",
	}
	entries := []DiffEntry{
		{Key: "URL", Left: "http://$BASE", Right: "http://$BASE", Status: StatusChanged},
	}

	result, err := ResolvedDiff(entries, leftEnv, rightEnv, ResolveOptions{Expand: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Left != "http://v1" {
		t.Errorf("expected Left='http://v1', got %q", result[0].Left)
	}
	if result[0].Right != "http://v2" {
		t.Errorf("expected Right='http://v2', got %q", result[0].Right)
	}
	if result[0].Status != StatusChanged {
		t.Errorf("expected status to be preserved as Changed")
	}
}

func TestResolvedDiff_NoExpand(t *testing.T) {
	leftEnv := map[string]string{"FOO": "$BAR"}
	rightEnv := map[string]string{"FOO": "$BAR"}
	entries := []DiffEntry{
		{Key: "FOO", Left: "$BAR", Right: "$BAR", Status: StatusUnchanged},
	}

	result, err := ResolvedDiff(entries, leftEnv, rightEnv, ResolveOptions{Expand: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Left != "$BAR" {
		t.Errorf("expected unexpanded value, got %q", result[0].Left)
	}
}

func TestResolvedDiff_StrictError(t *testing.T) {
	leftEnv := map[string]string{"FOO": "$MISSING"}
	rightEnv := map[string]string{"FOO": "bar"}
	entries := []DiffEntry{
		{Key: "FOO", Left: "$MISSING", Right: "bar", Status: StatusChanged},
	}

	_, err := ResolvedDiff(entries, leftEnv, rightEnv, ResolveOptions{Expand: true, Strict: true})
	if err == nil {
		t.Fatal("expected error in strict mode for missing variable")
	}
}

func TestResolvedDiff_PreservesOtherFields(t *testing.T) {
	leftEnv := map[string]string{"KEY": "left"}
	rightEnv := map[string]string{"KEY": "right"}
	entries := []DiffEntry{
		{Key: "KEY", Left: "left", Right: "right", Status: StatusChanged},
	}

	result, err := ResolvedDiff(entries, leftEnv, rightEnv, ResolveOptions{Expand: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Key != "KEY" {
		t.Errorf("expected Key='KEY', got %q", result[0].Key)
	}
	if result[0].Status != StatusChanged {
		t.Errorf("expected StatusChanged to be preserved")
	}
}
