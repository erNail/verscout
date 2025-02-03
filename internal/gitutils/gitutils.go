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

type TagInfo struct {
	Name     string
	UnixTime int64
	TagRef   *plumbing.Reference
}

func getCommitTimestamp(repo *git.Repository, hash plumbing.Hash) (int64, error) {
	commit, err := repo.CommitObject(hash)
	if err != nil {
		return 0, fmt.Errorf("failed to get commit object: %w", err)
	}

	return commit.Committer.When.Unix(), nil
}

func GetTagsWithTimestamps(repo *git.Repository) ([]TagInfo, error) {
	var tagsInfo []TagInfo

	tagRefs, err := repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tags: %w", err)
	}

	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		commit, err := repo.CommitObject(tagRef.Hash())
		if err != nil {
			return fmt.Errorf("failed to get commit object for tag %s: %w", tagRef.Name().Short(), err)
		}

		commitTime := commit.Committer.When.Unix()
		tagInfo := TagInfo{Name: tagRef.Name().Short(), UnixTime: commitTime, TagRef: tagRef}
		tagsInfo = append(tagsInfo, tagInfo)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate tags: %w", err)
	}

	return tagsInfo, nil
}

func GetLatestVersionTag(repo *git.Repository) (*TagInfo, error) {
	tags, err := GetTagsWithTimestamps(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags with timestamps: %w", err)
	}

	var latestTag TagInfo

	for _, tag := range tags {
		if semverutils.IsValidSemVerTag(tag.Name) {
			if latestTag.Name == "" || tag.UnixTime > latestTag.UnixTime {
				latestTag = tag
			}
		}
	}

	if latestTag.Name == "" {
		return nil, errors.New("no valid version tags found")
	}

	log.WithField("tag", latestTag.Name).Info("Found latest version tag")

	return &latestTag, nil
}

func GetLatestVersion(repo *git.Repository) (*semverutils.SemVer, error) {
	latestTag, err := GetLatestVersionTag(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version tag: %w", err)
	}

	semVer, err := semverutils.ExtractSemVerStruct(latestTag.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to extract semver struct: %w", err)
	}

	return semVer, nil
}

func GetCommitsSinceCommitHash(repo *git.Repository, commitHash plumbing.Hash) ([]*object.Commit, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()}) // Start from HEAD
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

	return commits, nil
}

func GetCommitMessagesSinceCommitHash(repo *git.Repository, commitHash plumbing.Hash) ([]string, error) {
	commits, err := GetCommitsSinceCommitHash(repo, commitHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits since commit hash: %w", err)
	}

	commitMessages := make([]string, 0, len(commits))
	for _, commit := range commits {
		commitMessages = append(commitMessages, commit.Message)
	}

	if len(commitMessages) == 0 {
		log.Info("No commits found since the given commit hash")
	} else {
		log.WithField("commitMessages", commitMessages).Info("Found commits since the given commit hash")
	}

	return commitMessages, nil
}
