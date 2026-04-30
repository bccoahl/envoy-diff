package dotenv_test

import (
	"testing"

	"github.com/envoy-diff/envoy-diff/internal/dotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func findEntry(entries []dotenv.DiffEntry, key string) (dotenv.DiffEntry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return dotenv.DiffEntry{}, false
}

func TestDiff_Added(t *testing.T) {
	a := dotenv.EnvMap{}
	b := dotenv.EnvMap{"NEW_KEY": "hello"}

	entries := dotenv.Diff(a, b)
	require.Len(t, entries, 1)
	assert.Equal(t, dotenv.Added, entries[0].Kind)
	assert.Equal(t, "NEW_KEY", entries[0].Key)
	assert.Equal(t, "hello", entries[0].ValueB)
}

func TestDiff_Removed(t *testing.T) {
	a := dotenv.EnvMap{"OLD_KEY": "bye"}
	b := dotenv.EnvMap{}

	entries := dotenv.Diff(a, b)
	require.Len(t, entries, 1)
	assert.Equal(t, dotenv.Removed, entries[0].Kind)
	assert.Equal(t, "bye", entries[0].ValueA)
}

func TestDiff_Changed(t *testing.T) {
	a := dotenv.EnvMap{"KEY": "old"}
	b := dotenv.EnvMap{"KEY": "new"}

	entries := dotenv.Diff(a, b)
	require.Len(t, entries, 1)
	assert.Equal(t, dotenv.Changed, entries[0].Kind)
	assert.Equal(t, "old", entries[0].ValueA)
	assert.Equal(t, "new", entries[0].ValueB)
}

func TestDiff_Unchanged(t *testing.T) {
	a := dotenv.EnvMap{"KEY": "same"}
	b := dotenv.EnvMap{"KEY": "same"}

	entries := dotenv.Diff(a, b)
	require.Len(t, entries, 1)
	assert.Equal(t, dotenv.Unchanged, entries[0].Kind)
}

func TestDiff_Empty(t *testing.T) {
	a := dotenv.EnvMap{}
	b := dotenv.EnvMap{}

	entries := dotenv.Diff(a, b)
	assert.Empty(t, entries)
}

func TestDiff_Mixed(t *testing.T) {
	a := dotenv.EnvMap{"KEEP": "v1", "REMOVE": "gone", "CHANGE": "old"}
	b := dotenv.EnvMap{"KEEP": "v1", "ADD": "new", "CHANGE": "new"}

	entries := dotenv.Diff(a, b)

	e, ok := findEntry(entries, "KEEP")
	require.True(t, ok)
	assert.Equal(t, dotenv.Unchanged, e.Kind)

	e, ok = findEntry(entries, "REMOVE")
	require.True(t, ok)
	assert.Equal(t, dotenv.Removed, e.Kind)

	e, ok = findEntry(entries, "ADD")
	require.True(t, ok)
	assert.Equal(t, dotenv.Added, e.Kind)

	e, ok = findEntry(entries, "CHANGE")
	require.True(t, ok)
	assert.Equal(t, dotenv.Changed, e.Kind)
}

func TestDiff_SortedOutput(t *testing.T) {
	a := dotenv.EnvMap{"ZZZ": "1", "AAA": "2", "MMM": "3"}
	b := dotenv.EnvMap{"ZZZ": "1", "AAA": "2", "MMM": "3"}

	entries := dotenv.Diff(a, b)
	require.Len(t, entries, 3)
	assert.Equal(t, "AAA", entries[0].Key)
	assert.Equal(t, "MMM", entries[1].Key)
	assert.Equal(t, "ZZZ", entries[2].Key)
}
