// Package dotenv provides functionality to load and parse .env files
// into a normalized key-value map for use in diffing.
package dotenv

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// EnvMap represents a set of environment variables as a key-value map.
type EnvMap map[string]string

// LoadFile reads a .env file from the given path and returns an EnvMap.
// Returns an error if the file does not exist or cannot be parsed.
func LoadFile(path string) (EnvMap, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("dotenv: file not found: %s", path)
	}

	raw, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("dotenv: failed to parse %s: %w", path, err)
	}

	return EnvMap(raw), nil
}

// LoadFiles reads multiple .env files and merges them into a single EnvMap.
// Later files take precedence over earlier ones on key conflicts.
func LoadFiles(paths ...string) (EnvMap, error) {
	merged := make(EnvMap)

	for _, p := range paths {
		em, err := LoadFile(p)
		if err != nil {
			return nil, err
		}
		for k, v := range em {
			merged[k] = v
		}
	}

	return merged, nil
}
