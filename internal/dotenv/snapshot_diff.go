package dotenv

import "fmt"

// SnapshotDiffResult holds the comparison between a live diff and a saved snapshot.
type SnapshotDiffResult struct {
	// New entries present in live but not in snapshot.
	New []DiffEntry
	// Gone entries present in snapshot but not in live.
	Gone []DiffEntry
	// Shifted entries whose status changed between snapshot and live.
	Shifted []ShiftedEntry
}

// ShiftedEntry describes a key whose DiffStatus changed since the snapshot.
type ShiftedEntry struct {
	Key         string
	OldStatus   DiffStatus
	NewStatus   DiffStatus
	CurrentEntry DiffEntry
}

// CompareSnapshot compares live diff entries against a previously saved snapshot.
// It returns a SnapshotDiffResult describing what changed since the snapshot was taken.
func CompareSnapshot(snap *Snapshot, live []DiffEntry) (*SnapshotDiffResult, error) {
	if snap == nil {
		return nil, fmt.Errorf("snapshot: nil snapshot provided")
	}

	snapMap := make(map[string]DiffEntry, len(snap.Entries))
	for _, e := range snap.Entries {
		snapMap[e.Key] = e
	}

	liveMap := make(map[string]DiffEntry, len(live))
	for _, e := range live {
		liveMap[e.Key] = e
	}

	result := &SnapshotDiffResult{}

	for _, le := range live {
		se, found := snapMap[le.Key]
		if !found {
			result.New = append(result.New, le)
			continue
		}
		if se.Status != le.Status {
			result.Shifted = append(result.Shifted, ShiftedEntry{
				Key:          le.Key,
				OldStatus:    se.Status,
				NewStatus:    le.Status,
				CurrentEntry: le,
			})
		}
	}

	for _, se := range snap.Entries {
		if _, found := liveMap[se.Key]; !found {
			result.Gone = append(result.Gone, se)
		}
	}

	return result, nil
}
