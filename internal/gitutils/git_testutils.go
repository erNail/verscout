//go:build test

package gitutils

import (
	"fmt"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

// Mock for GitRepository.
type MockGit struct {
	Repo *git.Repository
}

func (m *MockGit) PlainOpen(_ string) (*git.Repository, error) {
	return m.Repo, nil
}

func CreateTestRepo() (*git.Repository, error) {
	fs := memfs.New()

	repo, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	return repo, nil
}

func CreateTestCommit(repo *git.Repository, message, fileName, content string, time time.Time) (plumbing.Hash, error) {
	worktree, err := repo.Worktree()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to commit: %w", err)
	}

	file, err := worktree.Filesystem.Create(fileName)
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to create file: %w", err)
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to write file: %w", err)
	}

	file.Close()

	_, err = worktree.Add(fileName)
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to add file: %w", err)
	}

	commitHash, err := worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Author",
			Email: "author@test.com",
			When:  time,
		},
	})
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to commit: %w", err)
	}

	return commitHash, nil
}

func CreateTag(
	repo *git.Repository,
	tagName string,
	commitHash plumbing.Hash,
) (*plumbing.Reference, error) {
	tagHash, err := repo.CreateTag(tagName, commitHash, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create annotated tag: %w", err)
	}

	return tagHash, nil
}

func CreateAnnotatedTag(
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
