package atomicfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteFile_CreatesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.json")

	require.NoError(t, WriteFile(path, []byte("hello"), 0o600))

	data, err := os.ReadFile(path) // #nosec G304 // Test file path
	require.NoError(t, err)
	assert.Equal(t, "hello", string(data))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestWriteFile_ReplacesExisting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.json")
	require.NoError(t, os.WriteFile(path, []byte("old"), 0o600))

	require.NoError(t, WriteFile(path, []byte("new"), 0o600))

	data, err := os.ReadFile(path) // #nosec G304 // Test file path
	require.NoError(t, err)
	assert.Equal(t, "new", string(data))
}

func TestWriteFile_LeavesNoTempFileBehind(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")

	require.NoError(t, WriteFile(path, []byte("hello"), 0o600))

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "out.json", entries[0].Name())
}

func TestWriteFile_MissingDirectoryFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "out.json")

	err := WriteFile(path, []byte("hello"), 0o600)
	assert.Error(t, err)
}
