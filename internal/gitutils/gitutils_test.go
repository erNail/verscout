package gitutils

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestRepo() (*git.Repository, error) {
	fs := memfs.New()

	repo, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	return repo, nil
}

func createTestCommit(repo *git.Repository, message string, content string, time time.Time) (plumbing.Hash, error) {
	worktree, err := repo.Worktree()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to get worktree: %w", err)
	}

	file, err := worktree.Filesystem.Create("testfile.txt")
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to create file: %w", err)
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to write to file: %w", err)
	}

	file.Close()

	_, err = worktree.Add("testfile.txt")
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to add file to worktree: %w", err)
	}

	commitHash, err := worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Author",
			Email: "author@test.com",
			When:  time,
		},
	})
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to commit changes: %w", err)
	}

	return commitHash, nil
}

func createAnnotatedTag(
	repo *git.Repository,
	tagName string,
	commitHash plumbing.Hash,
	message string,
) (*plumbing.Reference, error) {
	tagHash, err := repo.CreateTag(tagName, commitHash, &git.CreateTagOptions{
		Message: message,
		Tagger: &object.Signature{
			Name:  "Test Tagger",
			Email: "tagger@test.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create annotated tag: %w", err)
	}

	return tagHash, nil
}

func TestGetTagsWithTimestamps_TwoTags(t *testing.T) {
	t.Parallel()

	repo, err := createTestRepo()
	require.NoError(t, err)
	commitHash, err := createTestCommit(repo, "First commit", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", commitHash, nil)
	require.NoError(t, err)
	commitHash, err = createTestCommit(repo, "Second commit", "Hello again, World!", time.Now().Add(1*time.Hour))
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.1", commitHash, nil)
	require.NoError(t, err)

	tagsInfos, err := GetTagsWithTimestamps(repo)
	require.NoError(t, err)
	assert.Len(t, tagsInfos, 2)
	assert.Equal(t, "1.0.0", tagsInfos[0].Name)
	assert.Equal(t, "1.0.1", tagsInfos[1].Name)
	assert.NotEqual(t, tagsInfos[0].UnixTime, tagsInfos[1].UnixTime)
}

func TestGetTagsWithTimestamps_NoTags(t *testing.T) {
	t.Parallel()

	repo, err := createTestRepo()
	require.NoError(t, err)

	tagsInfos, err := GetTagsWithTimestamps(repo)
	require.ErrorIs(t, err, ErrNoTags)
	assert.Empty(t, tagsInfos)
}

func TestGetTagsWithTimestamps_AnnotatedTag(t *testing.T) {
	t.Parallel()

	repo, err := createTestRepo()
	require.NoError(t, err)
	commitHash, err := createTestCommit(repo, "First commit", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = createAnnotatedTag(repo, "v1.0.0", commitHash, "Annotated tag")
	require.NoError(t, err)

	tagsInfos, err := GetTagsWithTimestamps(repo)
	require.NoError(t, err)
	assert.Len(t, tagsInfos, 1)
	assert.Equal(t, "v1.0.0", tagsInfos[0].Name)
	assert.NotZero(t, tagsInfos[0].UnixTime)
}

func TestGetCommitTimestamp(t *testing.T) {
	t.Parallel()

	repo, err := createTestRepo()
	require.NoError(t, err)
	commitHash, err := createTestCommit(repo, "First commit", "Hello, World!", time.Now())
	require.NoError(t, err)
	tagReference, err := repo.CreateTag("1.0.0", commitHash, nil)
	require.NoError(t, err)

	timestamp, err := getCommitTimestamp(repo, tagReference.Hash())
	require.NoError(t, err)
	assert.NotZero(t, timestamp)
}

func TestGetLatestVersion_NoTags(t *testing.T) {
	t.Parallel()

	repo, err := createTestRepo()
	require.NoError(t, err)

	tagInfo, err := GetLatestVersionTag(repo)
	require.ErrorIs(t, err, ErrNoTags)
	assert.Nil(t, tagInfo)
}

func TestGetLatestVersion_NoValidSemVerTags(t *testing.T) {
	t.Parallel()

	repo, err := createTestRepo()
	require.NoError(t, err)
	commitHash, err := createTestCommit(repo, "First commit", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.0.0", commitHash, nil)
	require.NoError(t, err)

	tagInfo, err := GetLatestVersionTag(repo)
	require.ErrorIs(t, err, ErrNoValidVersionTags)
	assert.Nil(t, tagInfo)
}

func TestGetLatestVersion_WithValidSemVerTags(t *testing.T) {
	t.Parallel()

	repo, err := createTestRepo()
	require.NoError(t, err)
	commitHash, err := createTestCommit(repo, "First commit", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", commitHash, nil)
	require.NoError(t, err)
	commitHash, err = createTestCommit(repo, "Second commit", "Hello again, World!", time.Now().Add(1*time.Hour))
	require.NoError(t, err)
	_, err = repo.CreateTag("1.1.0", commitHash, nil)
	require.NoError(t, err)
	commitHash, err = createTestCommit(repo, "Third commit", "Hello once more, World!", time.Now().Add(2*time.Hour))
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.1", commitHash, nil)
	require.NoError(t, err)

	tagInfo, err := GetLatestVersionTag(repo)
	require.NoError(t, err)
	assert.NotNil(t, tagInfo)
	assert.Equal(t, "1.0.1", tagInfo.Name)
}

func TestGetLatestVersion_WithMixedTags(t *testing.T) {
	t.Parallel()

	repo, err := createTestRepo()
	require.NoError(t, err)
	commitHash, err := createTestCommit(repo, "First commit", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", commitHash, nil)
	require.NoError(t, err)
	commitHash, err = createTestCommit(repo, "Second commit", "Hello again, World!", time.Now().Add(1*time.Hour))
	require.NoError(t, err)
	_, err = repo.CreateTag("not-a-semver", commitHash, nil)
	require.NoError(t, err)
	commitHash, err = createTestCommit(repo, "Third commit", "Hello once more, World!", time.Now().Add(2*time.Hour))
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.1", commitHash, nil)
	require.NoError(t, err)

	tagInfo, err := GetLatestVersionTag(repo)
	require.NoError(t, err)
	assert.NotNil(t, tagInfo)
	assert.Equal(t, "1.0.1", tagInfo.Name)
}

func TestGetCommitsSinceCommitHash(t *testing.T) {
	t.Parallel()

	repo, err := createTestRepo()
	require.NoError(t, err)

	commitHash1, err := createTestCommit(repo, "First commit", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", commitHash1, nil)
	require.NoError(t, err)

	commitHash2, err := createTestCommit(repo, "Second commit", "Hello again, World!", time.Now().Add(time.Hour))
	require.NoError(t, err)

	commitHash3, err := createTestCommit(repo, "Third commit", "Hello once more, World!", time.Now().Add(2*time.Hour))
	require.NoError(t, err)

	commits, err := GetCommitsSinceCommitHash(repo, commitHash1)
	require.NoError(t, err)
	assert.Len(t, commits, 2)
	assert.Equal(t, commitHash3, commits[0].Hash)
	assert.Equal(t, commitHash2, commits[1].Hash)
}

func TestGetCommitMessagesSinceCommitHash(t *testing.T) {
	t.Parallel()

	repo, err := createTestRepo()
	require.NoError(t, err)

	commitHash1, err := createTestCommit(repo, "First commit", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", commitHash1, nil)
	require.NoError(t, err)

	_, err = createTestCommit(repo, "Second commit", "Hello again, World!", time.Now().Add(1*time.Hour))
	require.NoError(t, err)

	_, err = createTestCommit(repo, "Third commit", "Hello once more, World!", time.Now().Add(2*time.Hour))
	require.NoError(t, err)

	messages, err := GetCommitMessagesSinceCommitHash(repo, commitHash1)
	require.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, "Third commit", messages[0])
	assert.Equal(t, "Second commit", messages[1])
}

func TestGetCommitFromTag_LightweightTag(t *testing.T) {
	t.Parallel()

	repo, err := createTestRepo()
	require.NoError(t, err)
	commitHash, err := createTestCommit(repo, "First commit", "Hello, World!", time.Now())
	require.NoError(t, err)
	tagRef, err := repo.CreateTag("v1.0.0", commitHash, nil)
	require.NoError(t, err)

	commit, err := GetCommitFromTag(repo, tagRef)
	require.NoError(t, err)
	assert.Equal(t, commitHash, commit.Hash)
}

func TestGetCommitFromTag_AnnotatedTag(t *testing.T) {
	t.Parallel()

	repo, err := createTestRepo()
	require.NoError(t, err)
	commitHash, err := createTestCommit(repo, "First commit", "Hello, World!", time.Now())
	require.NoError(t, err)
	tagRef, err := createAnnotatedTag(repo, "v1.0.0", commitHash, "Annotated tag")
	require.NoError(t, err)

	commit, err := GetCommitFromTag(repo, tagRef)
	require.NoError(t, err)
	assert.Equal(t, commitHash, commit.Hash)
}
