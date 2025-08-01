package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// GitInterface defines the required git operations for version management.
type GitInterface interface {
	PlainOpen(path string) (*git.Repository, error)
}

// Git implements the GitInterface for interacting with git repositories.
type Git struct{}

// PlainOpen opens a git repository at the specified path.
func (g *Git) PlainOpen(path string) (*git.Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository at %s: %w", path, err)
	}

	return repo, nil
}

// ExitError is a custom error type that includes an exit code and an underlying error.
// It is used to signal specific exit conditions for the CLI application.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	return e.Err.Error()
}

func (e *ExitError) Unwrap() error {
	return e.Err
}

// This will be set during the build via `-ldflags "-s -w -X github.com/erNail/verscout/cmd.version={{ .Version }}"`.
var version = "dev"

// NewRootCmd creates the root command for the CLI application.
// This command serves as the entry point and parent for all other commands.
func NewRootCmd() *cobra.Command {
	var repoDirectoryPath string

	rootCmd := &cobra.Command{
		Use:           "verscout",
		Short:         "Find the latest version tag and calculate the next version",
		Long:          `Find the latest version tag and calculate the next version based on conventional commits`,
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	rootCmd.PersistentFlags().StringVarP(&repoDirectoryPath, "dir", "d", ".", "directory path to the git repository")
	rootCmd.AddCommand(NewLatestCmd(&Git{}, &repoDirectoryPath))
	rootCmd.AddCommand(NewNextCmd(&Git{}, &repoDirectoryPath))

	return rootCmd
}

// Execute runs the root command.
func Execute() {
	cmd := NewRootCmd()

	err := cmd.Execute()
	if err != nil {
		var exitErr *ExitError
		if errors.As(err, &exitErr) {
			log.Error(err)
			os.Exit(exitErr.Code)
		}

		log.Fatal(err)
	}
}
