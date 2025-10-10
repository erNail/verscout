package gitutils

import (
	"errors"
	"fmt"

	"github.com/erNail/verscout/internal/semverutils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrNoCommitsFound indicates that no commits were found since the given hash.
	ErrNoCommitsFound = errors.New("no commits found since given hash")
	// ErrNoValidVersionTags indicates that no valid version tags were found.
	ErrNoValidVersionTags = errors.New("no valid version tags found")
	// ErrNoTags indicates that no tags were found in the repository.
	ErrNoTags = errors.New("no tags found")
)

// TagInfo holds information about a git tag.
type TagInfo struct {
	Name   string
	Commit *object.Commit
	TagRef *plumbing.Reference
}

// getCommitTimestamp retrieves the Unix timestamp of a commit given its hash.
func getCommitTimestamp(repo *git.Repository, hash plumbing.Hash) (int64, error) {
	commit, err := repo.CommitObject(hash)
	if err != nil {
		return 0, fmt.Errorf("failed to get commit object: %w", err)
	}

	return commit.Committer.When.Unix(), nil
}

// GetCommitFromTag retrieves the commit that a tag points to, handling both lightweight and annotated tags.
func GetCommitFromTag(repo *git.Repository, tagRef *plumbing.Reference) (*object.Commit, error) {
	tagObject, err := repo.TagObject(tagRef.Hash())
	if errors.Is(err, plumbing.ErrObjectNotFound) {
		// Lightweight tag, points directly to a commit
		commit, err := repo.CommitObject(tagRef.Hash())
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get commit object for lightweight tag %s: %w",
				tagRef.Name().Short(),
				err,
			)
		}

		return commit, nil
	}

	if err == nil {
		// Annotated tag, points to a tag object which points to a commit
		commit, err := repo.CommitObject(tagObject.Target)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get commit object for annotated tag %s: %w",
				tagRef.Name().Short(),
				err,
			)
		}

		return commit, nil
	}

	return nil, fmt.Errorf("failed to get tag object for tag %s: %w", tagRef.Name().Short(), err)
}

// GetTagsWithAssociatedCommits returns all tags in the repository with their associated commits.
// Returns ErrNoTags if no tags are found.
func GetTagsWithAssociatedCommits(repo *git.Repository) ([]TagInfo, error) {
	var tagsInfo []TagInfo

	tagRefs, err := repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tags: %w", err)
	}

	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		commit, err := GetCommitFromTag(repo, tagRef)
		if err != nil {
			return fmt.Errorf(
				"failed to get commit object for tag %s: %w",
				tagRef.Name().Short(),
				err,
			)
		}

		tagInfo := TagInfo{Name: tagRef.Name().Short(), Commit: commit, TagRef: tagRef}
		tagsInfo = append(tagsInfo, tagInfo)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate tags: %w", err)
	}

	if len(tagsInfo) == 0 {
		return nil, ErrNoTags
	}

	return tagsInfo, nil
}

// GetLatestVersionTag finds the most recent semantic version tag in the repository.
// Returns ErrNoValidVersionTags if no valid version tags are found.
// Returns ErrNoTags if no tags are found in the repository.
func GetLatestVersionTag(repo *git.Repository) (*TagInfo, error) {
	tags, err := GetTagsWithAssociatedCommits(repo)
	if err != nil {
		// Error type could be ErrNoTags
		return nil, fmt.Errorf("failed to get tags with timestamps: %w", err)
	}

	var latestTag TagInfo

	for _, tag := range tags {
		if semverutils.IsValidSemVerTag(tag.Name) {
			if latestTag.Name == "" || tag.Commit.Committer.When.Unix() > latestTag.Commit.Committer.When.Unix() {
				latestTag = tag
			}
		}
	}

	if latestTag.Name == "" {
		return nil, ErrNoValidVersionTags
	}

	log.WithField("tag", latestTag.Name).Info("Found latest version tag")

	return &latestTag, nil
}

// GetLatestVersion returns the latest semantic version as a SemVer struct.
// Returns ErrNoTags if no tags are found.
// Returns ErrNoValidVersionTags if no valid version tags are found.
func GetLatestVersion(repo *git.Repository) (*semverutils.SemVer, error) {
	latestTag, err := GetLatestVersionTag(repo)
	if err != nil {
		// Error type could be ErrNoValidVersionTags or ErrNoTags
		return nil, fmt.Errorf("failed to get latest version tag: %w", err)
	}

	semVer, err := semverutils.ExtractSemVerStruct(latestTag.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to extract semver struct: %w", err)
	}

	return semVer, nil
}

// GetCommitsSinceCommitHash returns all commits made after the specified commit hash.
// Returns ErrNoCommitsFound if no commits are found.
func GetCommitsSinceCommitHash(
	repo *git.Repository,
	commitHash plumbing.Hash,
) ([]*object.Commit, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commitIter, err := repo.Log(&git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	var commits []*object.Commit

	stopCollecting := false

	err = commitIter.ForEach(func(commit *object.Commit) error {
		if commit.Hash == commitHash {
			stopCollecting = true // Stop when we reach the given commit

			return nil
		}

		if !stopCollecting {
			commits = append(commits, commit)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	if len(commits) == 0 {
		return nil, ErrNoCommitsFound
	}

	return commits, nil
}

// GetCommitMessagesSinceCommitHash returns the commit messages for all commits made after the specified commit hash.
// Returns ErrNoCommitsFound if no commits are found.
func GetCommitMessagesSinceCommitHash(
	repo *git.Repository,
	commitHash plumbing.Hash,
) ([]string, error) {
	commits, err := GetCommitsSinceCommitHash(repo, commitHash)
	if err != nil {
		// Error type could be ErrNoCommitsFound
		return nil, fmt.Errorf("failed to get commits since commit hash: %w", err)
	}

	commitMessages := make([]string, 0, len(commits))

	for _, commit := range commits {
		log.WithField("commitMessage", commit.Message).Info("Found commit message")
		commitMessages = append(commitMessages, commit.Message)
	}

	return commitMessages, nil
}
