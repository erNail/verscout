package cmd

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
	TagName           string
	AdditionalCommits []string
}

func (m *MockGit) PlainOpen(_ string) (*git.Repository, error) {
	repository, err := createTestRepo()
	if err != nil {
		return nil, fmt.Errorf("failed to create test repo: %w", err)
	}

	commitHash, err := createTestCommit(repository, "Initial commit", "README.md", "Hello, World!", time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create test commit: %w", err)
	}

	if m.TagName != "" {
		_, err = repository.CreateTag(m.TagName, commitHash, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create tag: %w", err)
		}
	}

	for index, message := range m.AdditionalCommits {
		_, err = createTestCommit(
			repository,
			message,
			"README.md",
			"Hello again, World!",
			time.Now().Add(time.Duration(index+1)*time.Hour),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create additional commit: %w", err)
		}
	}

	return repository, nil
}

func createTestRepo() (*git.Repository, error) {
	fs := memfs.New()

	repo, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	return repo, nil
}

func createTestCommit(repo *git.Repository, message, fileName, content string, time time.Time) (plumbing.Hash, error) {
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
