package cmd

import (
	"bytes"
	"os"
	"path/filepath"
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

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0, ".verscout-config.yaml")
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

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0, ".verscout-config.yaml")
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

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0, ".verscout-config.yaml")
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

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0, ".verscout-config.yaml")
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

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0, ".verscout-config.yaml")
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

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0, ".verscout-config.yaml")
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

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0, ".verscout-config.yaml")
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

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0, ".verscout-config.yaml")
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

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 0, ".verscout-config.yaml")
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

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 2, ".verscout-config.yaml")
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

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 2, ".verscout-config.yaml")
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

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoDirectoryPath, 2, ".verscout-config.yaml")
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

func TestHandleNextCommand_CustomMajorBumpConfig(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)

	// Setup config file
	yamlContent := `
bumps:
  majorPatterns:
    - "(?m)^BREAK:"
  minorPatterns:
    - "^MINOR:"
  patchPatterns:
    - "^FIX:"
`
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(yamlContent), 0o600))

	// Setup repository
	commitHash, err := gitutils.CreateTestCommit(repo, "Initial commit", "test.txt", "test", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "v1.0.0", commitHash)
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(repo, "BREAK: this should trigger major bump", "test.txt", "test2", time.Now())
	require.NoError(t, err)

	var output bytes.Buffer

	repoPath := "."

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoPath, 0, configPath)
	require.NoError(t, err)
	assert.Equal(t, "2.0.0\n", output.String())
}

func TestHandleNextCommand_InvalidConfig(t *testing.T) {
	t.Parallel()

	repo, err := gitutils.CreateTestRepo()
	require.NoError(t, err)

	// Setup invalid config file
	content := []byte(`invalid: yaml: content`)
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(configPath, content, 0o600))

	// Setup repository
	commitHash, err := gitutils.CreateTestCommit(repo, "Initial commit", "test.txt", "test", time.Now())
	require.NoError(t, err)
	_, err = gitutils.CreateTag(repo, "v1.0.0", commitHash)
	require.NoError(t, err)
	_, err = gitutils.CreateTestCommit(repo, "fix: should use default config", "test.txt", "test2", time.Now())
	require.NoError(t, err)

	var output bytes.Buffer

	repoPath := "."

	err = HandleNextCommand(&output, &gitutils.MockGit{Repo: repo}, &repoPath, 0, configPath)
	require.NoError(t, err)
	assert.Equal(t, "1.0.1\n", output.String())
}
