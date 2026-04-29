package dotenv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var baseSortEntries = []DiffEntry{
	{Key: "ZEBRA", Status: StatusUnchanged, OldValue: "z", NewValue: "z"},
	{Key: "ALPHA", Status: StatusAdded, OldValue: "", NewValue: "1"},
	{Key: "MIDDLE", Status: StatusChanged, OldValue: "old", NewValue: "new"},
	{Key: "BRAVO", Status: StatusRemoved, OldValue: "2", NewValue: ""},
	{Key: "DELTA", Status: StatusAdded, OldValue: "", NewValue: "3"},
}

func TestSortEntries_ByKey(t *testing.T) {
	result := SortEntries(baseSortEntries, SortByKey)
	require.Len(t, result, 5)

	keys := make([]string, len(result))
	for i, e := range result {
		keys[i] = e.Key
	}
	assert.Equal(t, []string{"ALPHA", "BRAVO", "DELTA", "MIDDLE", "ZEBRA"}, keys)
}

func TestSortEntries_ByStatus(t *testing.T) {
	result := SortEntries(baseSortEntries, SortByStatus)
	require.Len(t, result, 5)

	// Added entries first (rank 0), then Removed (1), Changed (2), Unchanged (3)
	assert.Equal(t, StatusAdded, result[0].Status)
	assert.Equal(t, StatusAdded, result[1].Status)
	assert.Equal(t, StatusRemoved, result[2].Status)
	assert.Equal(t, StatusChanged, result[3].Status)
	assert.Equal(t, StatusUnchanged, result[4].Status)

	// Within the same status group entries should be sorted by key
	assert.Equal(t, "ALPHA", result[0].Key)
	assert.Equal(t, "DELTA", result[1].Key)
}

func TestSortEntries_None_PreservesOrder(t *testing.T) {
	result := SortEntries(baseSortEntries, SortNone)
	require.Len(t, result, 5)

	// Should be the same slice (no copy needed for SortNone)
	for i, e := range baseSortEntries {
		assert.Equal(t, e.Key, result[i].Key)
	}
}

func TestSortEntries_DoesNotMutateOriginal(t *testing.T) {
	original := make([]DiffEntry, len(baseSortEntries))
	copy(original, baseSortEntries)

	_ = SortEntries(baseSortEntries, SortByKey)

	for i, e := range baseSortEntries {
		assert.Equal(t, original[i].Key, e.Key, "original slice was mutated at index %d", i)
	}
}

func TestSortEntries_Empty(t *testing.T) {
	result := SortEntries([]DiffEntry{}, SortByKey)
	assert.Empty(t, result)
}
