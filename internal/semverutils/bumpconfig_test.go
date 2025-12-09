package semverutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadBumpConfigFromFile_Success(t *testing.T) {
	t.Parallel()

	yamlContent := `
bumps:
  majorPatterns:
    - "(?m)^MAJOR:"
    - "^BREAKING CHANGE:"
  minorPatterns:
    - "^MINOR:"
  patchPatterns:
    - "^FIX:"
`
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "bumpconfig.yaml")
	err := os.WriteFile(tmpFile, []byte(yamlContent), 0o600)
	require.NoError(t, err)

	config, err := LoadBumpConfigFromFile(tmpFile)
	require.NoError(t, err)

	assert.Len(t, config.Bumps.MajorPatterns, 2)
	assert.Equal(t, "(?m)^MAJOR:", config.Bumps.MajorPatterns[0])
	assert.Equal(t, "^BREAKING CHANGE:", config.Bumps.MajorPatterns[1])

	assert.Len(t, config.Bumps.MinorPatterns, 1)
	assert.Equal(t, "^MINOR:", config.Bumps.MinorPatterns[0])

	assert.Len(t, config.Bumps.PatchPatterns, 1)
	assert.Equal(t, "^FIX:", config.Bumps.PatchPatterns[0])
}

func TestLoadBumpConfigFromFile_FileNotFound(t *testing.T) {
	t.Parallel()

	_, err := LoadBumpConfigFromFile("nonexistent.yaml")
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestLoadBumpConfigFromFile_InvalidYAML(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid.yaml")
	err := os.WriteFile(tmpFile, []byte("not: [valid"), 0o600)
	require.NoError(t, err, "failed to write temp yaml file")

	_, err = LoadBumpConfigFromFile(tmpFile)
	require.Error(t, err)
	require.NotErrorIs(t, err, os.ErrNotExist)
}
