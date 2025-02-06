package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNextCmd_NoTags(t *testing.T) {
	t.Parallel()

	repoDirectoryPath := "."
	nextCmd := NewNextCmd(&MockGit{}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err := nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "1.0.0\n", output.String())
}

func TestNewNextCmd_ValidTag(t *testing.T) {
	t.Parallel()

	repoDirectoryPath := "."
	nextCmd := NewNextCmd(&MockGit{TagName: "1.0.0", AdditionalCommits: []string{"fix: bug fix"}}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err := nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "1.0.1\n", output.String())
}

func TestNewNextCmd_ValidTagWithVPrefix(t *testing.T) {
	t.Parallel()

	repoDirectoryPath := "."
	nextCmd := NewNextCmd(&MockGit{TagName: "v1.0.0", AdditionalCommits: []string{"fix: bug fix"}}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err := nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "1.0.1\n", output.String())
}

func TestNewNextCmd_InvalidTag(t *testing.T) {
	t.Parallel()

	repoDirectoryPath := "."
	nextCmd := NewNextCmd(&MockGit{
		TagName:           "invalid-tag",
		AdditionalCommits: []string{"fix: bug fix"},
	}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err := nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "1.0.0\n", output.String())
}

func TestNewNextCmd_MajorBump(t *testing.T) {
	t.Parallel()

	repoDirectoryPath := "."
	nextCmd := NewNextCmd(&MockGit{
		TagName: "1.0.0",
		AdditionalCommits: []string{
			"feat: NEW FEATURE\nBREAKING CHANGE: this is a breaking change",
		},
	}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err := nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "2.0.0\n", output.String())
}

func TestNewNextCmd_MinorBump(t *testing.T) {
	t.Parallel()

	repoDirectoryPath := "."
	nextCmd := NewNextCmd(&MockGit{TagName: "1.0.0", AdditionalCommits: []string{"feat: new feature"}}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err := nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "1.1.0\n", output.String())
}

func TestNewNextCmd_NoBumpChore(t *testing.T) {
	t.Parallel()

	repoDirectoryPath := "."
	nextCmd := NewNextCmd(&MockGit{
		TagName:           "1.0.0",
		AdditionalCommits: []string{"chore: update readme"},
	}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err := nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "", output.String())
}

func TestNewNextCmd_NoBumpNoAdditionalCommits(t *testing.T) {
	t.Parallel()

	repoDirectoryPath := "."
	nextCmd := NewNextCmd(&MockGit{TagName: "1.0.0"}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err := nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "", output.String())
}
