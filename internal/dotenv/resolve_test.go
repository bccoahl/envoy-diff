package dotenv

import (
	"testing"
)

func TestResolve_NoExpand(t *testing.T) {
	env := map[string]string{
		"FOO": "$BAR",
		"BAR": "hello",
	}
	result, err := Resolve(env, ResolveOptions{Expand: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "$BAR" {
		t.Errorf("expected FOO=$BAR, got %q", result["FOO"])
	}
}

func TestResolve_ExpandSimple(t *testing.T) {
	env := map[string]string{
		"BASE": "world",
		"GREETING": "hello $BASE",
	}
	result, err := Resolve(env, ResolveOptions{Expand: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["GREETING"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", result["GREETING"])
	}
}

func TestResolve_ExpandBraces(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"URL":  "http://${HOST}:8080",
	}
	result, err := Resolve(env, ResolveOptions{Expand: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["URL"] != "http://localhost:8080" {
		t.Errorf("expected 'http://localhost:8080', got %q", result["URL"])
	}
}

func TestResolve_StrictMissingVar(t *testing.T) {
	env := map[string]string{
		"FOO": "$MISSING",
	}
	_, err := Resolve(env, ResolveOptions{Expand: true, Strict: true})
	if err == nil {
		t.Fatal("expected error for missing variable in strict mode")
	}
}

func TestResolve_NonStrictMissingVar(t *testing.T) {
	env := map[string]string{
		"FOO": "$MISSING",
	}
	result, err := Resolve(env, ResolveOptions{Expand: true, Strict: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "" {
		t.Errorf("expected empty string for missing var, got %q", result["FOO"])
	}
}

func TestResolve_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{
		"BASE": "base",
		"FULL": "$BASE-extra",
	}
	orig := map[string]string{
		"BASE": "base",
		"FULL": "$BASE-extra",
	}
	_, _ = Resolve(env, ResolveOptions{Expand: true})
	for k, v := range orig {
		if env[k] != v {
			t.Errorf("input mutated: key %q changed to %q", k, env[k])
		}
	}
}
