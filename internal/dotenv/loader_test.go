package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-diff/envoy-diff/internal/dotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	err := os.WriteFile(p, []byte(content), 0600)
	require.NoError(t, err)
	return p
}

func TestLoadFile_Basic(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	em, err := dotenv.LoadFile(p)
	require.NoError(t, err)
	assert.Equal(t, "bar", em["FOO"])
	assert.Equal(t, "qux", em["BAZ"])
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := dotenv.LoadFile("/nonexistent/path/.env")
	assert.ErrorContains(t, err, "file not found")
}

func TestLoadFile_Empty(t *testing.T) {
	p := writeTempEnv(t, "")
	em, err := dotenv.LoadFile(p)
	require.NoError(t, err)
	assert.Empty(t, em)
}

func TestLoadFiles_Merge(t *testing.T) {
	p1 := writeTempEnv(t, "FOO=first\nSHARED=from_first\n")
	p2 := writeTempEnv(t, "BAR=second\nSHARED=from_second\n")

	em, err := dotenv.LoadFiles(p1, p2)
	require.NoError(t, err)
	assert.Equal(t, "first", em["FOO"])
	assert.Equal(t, "second", em["BAR"])
	// second file wins on conflict
	assert.Equal(t, "from_second", em["SHARED"])
}

func TestLoadFiles_ErrorPropagates(t *testing.T) {
	p := writeTempEnv(t, "KEY=value\n")
	_, err := dotenv.LoadFiles(p, "/no/such/file")
	assert.Error(t, err)
}
