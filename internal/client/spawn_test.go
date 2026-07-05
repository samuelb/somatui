package client

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenServerLog_CreatesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "server.log")

	f, err := openServerLog(path)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	assert.FileExists(t, path)
}

func TestOpenServerLog_AppendsBelowCap(t *testing.T) {
	path := filepath.Join(t.TempDir(), "server.log")
	require.NoError(t, os.WriteFile(path, []byte("old\n"), 0o600))

	f, err := openServerLog(path)
	require.NoError(t, err)
	_, err = f.WriteString("new\n")
	require.NoError(t, err)
	require.NoError(t, f.Close())

	data, err := os.ReadFile(path) // #nosec G304 // Test file path
	require.NoError(t, err)
	assert.Equal(t, "old\nnew\n", string(data))
}

func TestOpenServerLog_TruncatesAboveCap(t *testing.T) {
	path := filepath.Join(t.TempDir(), "server.log")
	oversized := bytes.Repeat([]byte("x"), maxServerLogSize+1)
	require.NoError(t, os.WriteFile(path, oversized, 0o600))

	f, err := openServerLog(path)
	require.NoError(t, err)
	_, err = f.WriteString("fresh\n")
	require.NoError(t, err)
	require.NoError(t, f.Close())

	data, err := os.ReadFile(path) // #nosec G304 // Test file path
	require.NoError(t, err)
	assert.Equal(t, "fresh\n", string(data), "an oversized log must be truncated at spawn")
}
