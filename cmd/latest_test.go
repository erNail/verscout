package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLatestCmd_ValidTag(t *testing.T) {
	t.Parallel()

	repoDirectoryPath := "."
	latestCmd := NewLatestCmd(&MockGit{TagName: "1.0.0"}, &repoDirectoryPath)

	var output bytes.Buffer

	latestCmd.SetOut(&output)
	err := latestCmd.Execute()
	require.NoError(t, err)

	assert.Contains(t, output.String(), "1.0.0")
}

func TestNewLatestCmd_ValidTagWithVPrefix(t *testing.T) {
	t.Parallel()

	repoDirectoryPath := "."
	latestCmd := NewLatestCmd(&MockGit{TagName: "v1.0.0"}, &repoDirectoryPath)

	var output bytes.Buffer

	latestCmd.SetOut(&output)
	err := latestCmd.Execute()
	require.NoError(t, err)

	assert.Contains(t, output.String(), "1.0.0")
	assert.NotContains(t, output.String(), "v1.0.0")
}

func TestNewLatestCmd_InvalidTag(t *testing.T) {
	t.Parallel()

	repoDirectoryPath := "."
	latestCmd := NewLatestCmd(&MockGit{TagName: "invalid-tag"}, &repoDirectoryPath)

	var output bytes.Buffer

	latestCmd.SetOut(&output)
	err := latestCmd.Execute()
	require.NoError(t, err)

	assert.Contains(t, output.String(), "")
}

func TestNewLatestCmd_NoTag(t *testing.T) {
	t.Parallel()

	repoDirectoryPath := "."
	latestCmd := NewLatestCmd(&MockGit{}, &repoDirectoryPath)

	var output bytes.Buffer

	latestCmd.SetOut(&output)
	err := latestCmd.Execute()
	require.NoError(t, err)

	assert.Contains(t, output.String(), "")
}
