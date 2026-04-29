package dotenv

// ResolvedDiff runs Resolve on both sides of a diff and returns new DiffEntries
// with the expanded values. The Status field is preserved from the original.
func ResolvedDiff(entries []DiffEntry, leftEnv, rightEnv map[string]string, opts ResolveOptions) ([]DiffEntry, error) {
	resolvedLeft, err := Resolve(leftEnv, opts)
	if err != nil {
		return nil, err
	}
	resolvedRight, err := Resolve(rightEnv, opts)
	if err != nil {
		return nil, err
	}

	result := make([]DiffEntry, len(entries))
	for i, e := range entries {
		updated := e
		if v, ok := resolvedLeft[e.Key]; ok {
			updated.Left = v
		}
		if v, ok := resolvedRight[e.Key]; ok {
			updated.Right = v
		}
		result[i] = updated
	}
	return result, nil
}
