package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/erNail/verscout/internal/gitutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNextCmd_NoExistingTags(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)

	repoDirectoryPath := "."
	nextCmd := NewNextCmd(&gitutils.MockGit{Repo: repo}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err = nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "1.0.0\n", output.String())
}

func TestNewNextCmd_ValidExistingTag(t *testing.T) {
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
	nextCmd := NewNextCmd(&gitutils.MockGit{Repo: repo}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err = nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "1.0.1\n", output.String())
}

func TestNewNextCmd_ValidExistingTagWithVPrefix(t *testing.T) {
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
	nextCmd := NewNextCmd(&gitutils.MockGit{Repo: repo}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err = nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "1.0.1\n", output.String())
}

func TestNewNextCmd_InvalidExistingTag(t *testing.T) {
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
	nextCmd := NewNextCmd(&gitutils.MockGit{Repo: repo}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err = nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "1.0.0\n", output.String())
}

func TestNewNextCmd_MajorBump(t *testing.T) {
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
	nextCmd := NewNextCmd(&gitutils.MockGit{Repo: repo}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err = nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "2.0.0\n", output.String())
}

func TestNewNextCmd_MinorBump(t *testing.T) {
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
	nextCmd := NewNextCmd(&gitutils.MockGit{Repo: repo}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err = nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "1.1.0\n", output.String())
}

func TestNewNextCmd_NoBumpChore(t *testing.T) {
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
	nextCmd := NewNextCmd(&gitutils.MockGit{Repo: repo}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err = nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "", output.String())
}

func TestNewNextCmd_NoBumpNoAdditionalCommits(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := gitutils.CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "1.0.0", commitHash)
	require.NoError(t, err)

	repoDirectoryPath := "."
	nextCmd := NewNextCmd(&gitutils.MockGit{Repo: repo}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err = nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "", output.String())
}

func TestNewNextCmd_AnnotatedExistingTag(t *testing.T) {
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
	nextCmd := NewNextCmd(&gitutils.MockGit{Repo: repo}, &repoDirectoryPath)

	var output bytes.Buffer

	nextCmd.SetOut(&output)
	err = nextCmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "1.0.1\n", output.String())
}
