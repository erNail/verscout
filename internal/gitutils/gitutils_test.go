package gitutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTagsWithTimestamps_TwoTags(t *testing.T) {
	t.Parallel()

	repo, err := CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", commitHash, nil)
	require.NoError(t, err)
	commitHash, err = CreateTestCommit(
		repo,
		"Second commit",
		"README.md",
		"Hello again, World!",
		time.Now().Add(1*time.Hour),
	)
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

	repo, err := CreateTestRepo()
	require.NoError(t, err)

	tagsInfos, err := GetTagsWithTimestamps(repo)
	require.ErrorIs(t, err, ErrNoTags)
	assert.Empty(t, tagsInfos)
}

func TestGetTagsWithTimestamps_AnnotatedTag(t *testing.T) {
	t.Parallel()

	repo, err := CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = CreateAnnotatedTag(repo, "v1.0.0", commitHash, "Annotated tag")
	require.NoError(t, err)

	tagsInfos, err := GetTagsWithTimestamps(repo)
	require.NoError(t, err)
	assert.Len(t, tagsInfos, 1)
	assert.Equal(t, "v1.0.0", tagsInfos[0].Name)
	assert.NotZero(t, tagsInfos[0].UnixTime)
}

func TestGetTagsWithTimestamps_LightweightTag(t *testing.T) {
	t.Parallel()

	repo, err := CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = CreateTag(repo, "v1.0.0", commitHash)
	require.NoError(t, err)

	tagsInfos, err := GetTagsWithTimestamps(repo)
	require.NoError(t, err)
	assert.Len(t, tagsInfos, 1)
	assert.Equal(t, "v1.0.0", tagsInfos[0].Name)
	assert.NotZero(t, tagsInfos[0].UnixTime)
}

func TestGetCommitTimestamp(t *testing.T) {
	t.Parallel()

	repo, err := CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	tagReference, err := repo.CreateTag("1.0.0", commitHash, nil)
	require.NoError(t, err)

	timestamp, err := getCommitTimestamp(repo, tagReference.Hash())
	require.NoError(t, err)
	assert.NotZero(t, timestamp)
}

func TestGetLatestVersion_NoTags(t *testing.T) {
	t.Parallel()

	repo, err := CreateTestRepo()
	require.NoError(t, err)

	tagInfo, err := GetLatestVersionTag(repo)
	require.ErrorIs(t, err, ErrNoTags)
	assert.Nil(t, tagInfo)
}

func TestGetLatestVersion_NoValidSemVerTags(t *testing.T) {
	t.Parallel()

	repo, err := CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.0.0", commitHash, nil)
	require.NoError(t, err)

	tagInfo, err := GetLatestVersionTag(repo)
	require.ErrorIs(t, err, ErrNoValidVersionTags)
	assert.Nil(t, tagInfo)
}

func TestGetLatestVersion_WithValidSemVerTags(t *testing.T) {
	t.Parallel()

	repo, err := CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", commitHash, nil)
	require.NoError(t, err)
	commitHash, err = CreateTestCommit(
		repo,
		"Second commit",
		"README.md",
		"Hello again, World!",
		time.Now().Add(1*time.Hour),
	)
	require.NoError(t, err)
	_, err = repo.CreateTag("1.1.0", commitHash, nil)
	require.NoError(t, err)
	commitHash, err = CreateTestCommit(
		repo,
		"Third commit",
		"README.md",
		"Hello once more, World!",
		time.Now().Add(2*time.Hour),
	)
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

	repo, err := CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", commitHash, nil)
	require.NoError(t, err)
	commitHash, err = CreateTestCommit(
		repo,
		"Second commit",
		"README.md",
		"Hello again, World!",
		time.Now().Add(1*time.Hour),
	)
	require.NoError(t, err)
	_, err = repo.CreateTag("not-a-semver", commitHash, nil)
	require.NoError(t, err)
	commitHash, err = CreateTestCommit(
		repo,
		"Third commit",
		"README.md",
		"Hello once more, World!",
		time.Now().Add(2*time.Hour),
	)
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

	repo, err := CreateTestRepo()
	require.NoError(t, err)

	commitHash1, err := CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", commitHash1, nil)
	require.NoError(t, err)

	commitHash2, err := CreateTestCommit(
		repo,
		"Second commit",
		"README.md",
		"Hello again, World!",
		time.Now().Add(time.Hour),
	)
	require.NoError(t, err)

	commitHash3, err := CreateTestCommit(
		repo,
		"Third commit",
		"README.md",
		"Hello once more, World!",
		time.Now().Add(2*time.Hour),
	)
	require.NoError(t, err)

	commits, err := GetCommitsSinceCommitHash(repo, commitHash1)
	require.NoError(t, err)
	assert.Len(t, commits, 2)
	assert.Equal(t, commitHash3, commits[0].Hash)
	assert.Equal(t, commitHash2, commits[1].Hash)
}

func TestGetCommitMessagesSinceCommitHash(t *testing.T) {
	t.Parallel()

	repo, err := CreateTestRepo()
	require.NoError(t, err)

	commitHash1, err := CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", commitHash1, nil)
	require.NoError(t, err)

	_, err = CreateTestCommit(repo, "Second commit", "README.md", "Hello again, World!", time.Now().Add(1*time.Hour))
	require.NoError(t, err)

	_, err = CreateTestCommit(repo, "Third commit", "README.md", "Hello once more, World!", time.Now().Add(2*time.Hour))
	require.NoError(t, err)

	messages, err := GetCommitMessagesSinceCommitHash(repo, commitHash1)
	require.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, "Third commit", messages[0])
	assert.Equal(t, "Second commit", messages[1])
}

func TestGetCommitFromTag_LightweightTag(t *testing.T) {
	t.Parallel()

	repo, err := CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	tagRef, err := repo.CreateTag("v1.0.0", commitHash, nil)
	require.NoError(t, err)

	commit, err := GetCommitFromTag(repo, tagRef)
	require.NoError(t, err)
	assert.Equal(t, commitHash, commit.Hash)
}

func TestGetCommitFromTag_AnnotatedTag(t *testing.T) {
	t.Parallel()

	repo, err := CreateTestRepo()
	require.NoError(t, err)
	commitHash, err := CreateTestCommit(repo, "First commit", "README.md", "Hello, World!", time.Now())
	require.NoError(t, err)
	tagRef, err := CreateAnnotatedTag(repo, "v1.0.0", commitHash, "Annotated tag")
	require.NoError(t, err)

	commit, err := GetCommitFromTag(repo, tagRef)
	require.NoError(t, err)
	assert.Equal(t, commitHash, commit.Hash)
}
