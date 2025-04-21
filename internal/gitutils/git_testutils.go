//go:build test

// Package gitutils provides testing utilities for working with Git repositories.
// It includes mock implementations and helper functions for creating test repositories,
// commits, and tags in memory.
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

// MockGit provides a mock implementation of the GitRepository interface for testing.
type MockGit struct {
	Repo *git.Repository
}

// PlainOpen implements the GitRepository interface by returning the mock repository.
// The path parameter is ignored in this mock implementation.
func (m *MockGit) PlainOpen(_ string) (*git.Repository, error) {
	return m.Repo, nil
}

// CreateTestRepo creates a new in-memory Git repository for testing purposes.
func CreateTestRepo() (*git.Repository, error) {
	fs := memfs.New()

	repo, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	return repo, nil
}

// CreateTestCommit creates a new commit in the given repository with the specified message,
// file name, content, and timestamp.
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

	err = file.Close()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to close file: %w", err)
	}

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

// CreateTag creates a lightweight tag in the repository pointing to the given commit.
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

// CreateAnnotatedTag creates a new annotated tag in the repository pointing to the given commit
// with the specified message and default tagger information.
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
