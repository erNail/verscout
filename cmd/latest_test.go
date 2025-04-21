package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/erNail/verscout/internal/gitutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleLatestCommand_ValidTag(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "1.0.0", commitHash)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleLatestCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Equal(t, "1.0.0\n", output.String())
}

func TestHandleLatestCommand_ValidTagWithVPrefix(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "v1.0.0", commitHash)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleLatestCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Equal(t, "1.0.0\n", output.String())
}

func TestHandleLatestCommand_InvalidTag(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "invalid-tag", commitHash)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleLatestCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Empty(t, output.String())
}

func TestHandleLatestCommand_NoTag(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleLatestCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Empty(t, output.String())
}

func TestHandleLatestCommand_AnnotatedTag(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateAnnotatedTag(repo, "1.0.0", commitHash, "Annotated tag")
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleLatestCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Equal(t, "1.0.0\n", output.String())
}

func TestHandleLatestCommand_NoLatestVersionExitCode_NoExistingTags(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleLatestCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 2)
	require.Error(t, err)

	var exitErr *ExitError

	require.ErrorAs(t, err, &exitErr)
	assert.Equal(t, 2, exitErr.Code)
}

func TestHandleLatestCommand_NoLatestVersionExitCode_NoValidVersion(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateAnnotatedTag(repo, "invalid-tag", commitHash, "Annotated tag")
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleLatestCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 2)
	require.Error(t, err)

	var exitErr *ExitError

	require.ErrorAs(t, err, &exitErr)
	assert.Equal(t, 2, exitErr.Code)
}

func TestHandleLatestCommand_NoLatestVersionExitCode_ExistingTag(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateAnnotatedTag(repo, "1.0.0", commitHash, "Annotated tag")
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleLatestCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 2)
	require.NoError(t, err)

	assert.Equal(t, "1.0.0\n", output.String())
}

func TestNewLatestCommand(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)

	repoDirectoryPath := "."

	cmd := NewLatestCmd(&gitutils.MockGit{Repo: repo}, &repoDirectoryPath)
	err = cmd.Execute()
	require.NoError(t, err)
}
