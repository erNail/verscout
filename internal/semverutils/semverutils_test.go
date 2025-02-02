package semverutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsValidSemVer_Valid(t *testing.T) {
	t.Parallel()
	assert.True(t, IsValidSemVer("1.0.0"))
	assert.True(t, IsValidSemVer("0.1.0"))
	assert.True(t, IsValidSemVer("0.0.1"))
	assert.True(t, IsValidSemVer("1.2.3"))
}

func TestIsValidSemVer_Invalid(t *testing.T) {
	t.Parallel()
	assert.False(t, IsValidSemVer("1.0"))
	assert.False(t, IsValidSemVer("1.0.0.0"))
	assert.False(t, IsValidSemVer("1.0.a"))
	assert.False(t, IsValidSemVer("a.b.c"))
	assert.False(t, IsValidSemVer("v1.2.3"))
}

func TestExtractSemVerStruct_Valid(t *testing.T) {
	t.Parallel()

	semVer, err := ExtractSemVerStruct("1.2.3")
	require.NoError(t, err)
	assert.Equal(t, &SemVer{Major: 1, Minor: 2, Patch: 3}, semVer)
}

func TestExtractSemVerStruct_InvalidPrefix(t *testing.T) {
	t.Parallel()

	semVer, err := ExtractSemVerStruct("v1.2.3")
	require.Error(t, err)
	assert.Nil(t, semVer)
}

func TestExtractSemVerStruct_InvalidFormat(t *testing.T) {
	t.Parallel()

	semVer, err := ExtractSemVerStruct("1.0")
	require.Error(t, err)
	assert.Nil(t, semVer)
}

func TestSemVerString(t *testing.T) {
	t.Parallel()

	semVer := &SemVer{Major: 1, Minor: 2, Patch: 3}
	assert.Equal(t, "1.2.3", semVer.String())

	semVer = &SemVer{Major: 0, Minor: 1, Patch: 0}
	assert.Equal(t, "0.1.0", semVer.String())

	semVer = &SemVer{Major: 2, Minor: 0, Patch: 1}
	assert.Equal(t, "2.0.1", semVer.String())
}

func TestCalculateNextVersion_BugFix(t *testing.T) {
	t.Parallel()

	nextVersion, err := CalculateNextVersion("1.0.0", []string{"fix: bug fix"})
	require.NoError(t, err)
	assert.Equal(t, "1.0.1", nextVersion)
}

func TestCalculateNextVersion_NewFeature(t *testing.T) {
	t.Parallel()

	nextVersion, err := CalculateNextVersion("1.0.0", []string{"feat: new feature"})
	require.NoError(t, err)
	assert.Equal(t, "1.1.0", nextVersion)
}

func TestCalculateNextVersion_BugFixAndNewFeature(t *testing.T) {
	t.Parallel()

	nextVersion, err := CalculateNextVersion("1.0.0", []string{"fix: bug fix", "feat: new feature"})
	require.NoError(t, err)
	assert.Equal(t, "1.1.0", nextVersion)
}

func TestCalculateNextVersion_BugFixAndBreakingChange(t *testing.T) {
	t.Parallel()

	nextVersion, err := CalculateNextVersion("1.0.0", []string{"fix: bug fix\n\nBREAKING CHANGE: major update"})
	require.NoError(t, err)
	assert.Equal(t, "2.0.0", nextVersion)
}

func TestCalculateNextVersion_NewFeatureAndBreakingChange(t *testing.T) {
	t.Parallel()

	nextVersion, err := CalculateNextVersion("1.0.0", []string{"feat: new feature\n\nBREAKING CHANGE: major update"})
	require.NoError(t, err)
	assert.Equal(t, "2.0.0", nextVersion)
}

func TestCalculateNextVersion_BugFixNewFeatureAndBreakingChange(t *testing.T) {
	t.Parallel()

	nextVersion, err := CalculateNextVersion("1.0.0", []string{
		"fix: bug fix",
		"feat: new feature\n\nBREAKING CHANGE: major update",
	})
	require.NoError(t, err)
	assert.Equal(t, "2.0.0", nextVersion)
}

func TestCalculateNextVersion_ChoreCommit(t *testing.T) {
	t.Parallel()

	nextVersion, err := CalculateNextVersion("1.0.0", []string{"chore: update readme"})
	require.Error(t, err)
	assert.Empty(t, nextVersion)
}

func TestCalculateNextVersion_NoCommits(t *testing.T) {
	t.Parallel()

	nextVersion, err := CalculateNextVersion("1.0.0", []string{})
	require.Error(t, err)
	assert.Empty(t, nextVersion)
}
