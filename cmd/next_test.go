package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/erNail/verscout/internal/gitutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleNextCommand_NoExistingTags(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Equal(t, "1.0.0\n", output.String())
}

func TestHandleNextCommand_ValidExistingTag(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "1.0.0", commitHash)
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(
		repo,
		"fix: Second commit",
		"README.md",
		"Hello, World! Again!",
		time.Now(),
	)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Equal(t, "1.0.1\n", output.String())
}

func TestHandleNextCommand_ValidExistingTagWithVPrefix(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "v1.0.0", commitHash)
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(
		repo,
		"fix: Second commit",
		"README.md",
		"Hello, World! Again!",
		time.Now(),
	)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Equal(t, "1.0.1\n", output.String())
}

func TestHandleNextCommand_InvalidExistingTag(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "invalid-tag", commitHash)
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(
		repo,
		"fix: Second commit",
		"README.md",
		"Hello, World! Again!",
		time.Now(),
	)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Equal(t, "1.0.0\n", output.String())
}

func TestHandleNextCommand_MajorBump(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "1.0.0", commitHash)
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(
		repo,
		"fix: Second commit\nBREAKING CHANGE: Break",
		"README.md",
		"Hello, World! Again!",
		time.Now(),
	)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Equal(t, "2.0.0\n", output.String())
}

func TestHandleNextCommand_MinorBump(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "1.0.0", commitHash)
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(
		repo,
		"feat: Second commit",
		"README.md",
		"Hello, World! Again!",
		time.Now(),
	)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Equal(t, "1.1.0\n", output.String())
}

func TestHandleNextCommand_NoBumpChore(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "1.0.0", commitHash)
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(
		repo,
		"chore: Second commit",
		"README.md",
		"Hello, World! Again!",
		time.Now(),
	)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Empty(t, output.String())
}

func TestHandleNextCommand_NoBumpNoAdditionalCommits(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "1.0.0", commitHash)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Empty(t, output.String())
}

func TestHandleNextCommand_AnnotatedExistingTag(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateAnnotatedTag(repo, "1.0.0", commitHash, "Annotated tag")
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(
		repo,
		"fix: Second commit",
		"README.md",
		"Hello, World! Again!",
		time.Now(),
	)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0)
	require.NoError(t, err)

	assert.Equal(t, "1.0.1\n", output.String())
}

func TestHandleNextCommand_NoNextVersionExitCode_NoNewCommits(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateAnnotatedTag(repo, "1.0.0", commitHash, "Annotated tag")
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 2)
	require.Error(t, err)

	var exitErr *ExitError

	require.ErrorAs(t, err, &exitErr)
	assert.Equal(t, 2, exitErr.Code)
}

func TestHandleNextCommand_NoNextVersionExitCode_NoBumpCommits(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateAnnotatedTag(repo, "1.0.0", commitHash, "Annotated tag")
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(
		repo,
		"chore: Second commit",
		"README.md",
		"Hello, World! Again!",
		time.Now(),
	)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 2)
	require.Error(t, err)

	var exitErr *ExitError

	require.ErrorAs(t, err, &exitErr)
	assert.Equal(t, 2, exitErr.Code)
}

func TestHandleNextCommand_NoNextVersionExitCode_FixCommit(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateAnnotatedTag(repo, "1.0.0", commitHash, "Annotated tag")
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(
		repo,
		"fix: Second commit",
		"README.md",
		"Hello, World! Again!",
		time.Now(),
	)
	require.NoError(t, err)

	repoDirectoryPath := "."

	var output bytes.Buffer

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 2)
	require.NoError(t, err)

	assert.Equal(t, "1.0.1\n", output.String())
}

func TestNewNextCommand(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)

	repoDirectoryPath := "."

	cmd := NewNextCmd(&gitutils.MockGit{Repo: repo}, &repoDirectoryPath)
	err = cmd.Execute()
	require.NoError(t, err)
}
