package cmd

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type GitInterface interface {
	PlainOpen(path string) (*git.Repository, error)
}

type Git struct{}

func (g *Git) PlainOpen(path string) (*git.Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository at %s: %w", path, err)
	}

	return repo, nil
}

// This will be set during the build via `-ldflags "-s -w -X github.com/erNail/verscout/cmd.version={{ .Version }}"`.
var version = "dev"

// NewRootCmd creates the root command for the CLI application.
// This command serves as the entry point and parent for all other commands.
//
// Returns:
//   - *cobra.Command: A pointer to the newly created cobra.Command.
func NewRootCmd() *cobra.Command {
	var repoDirectoryPath string

	rootCmd := &cobra.Command{
		Use:     "verscout",
		Short:   "Find the latest version tag and calculate the next version",
		Long:    `Find the latest version tag and calculate the next version based on conventional commits`,
		Version: version,
	}

	rootCmd.PersistentFlags().StringVarP(&repoDirectoryPath, "dir", "d", ".", "directory path to the git repository")
	rootCmd.AddCommand(NewLatestCmd(&Git{}, &repoDirectoryPath))
	rootCmd.AddCommand(NewNextCmd(&Git{}, &repoDirectoryPath))

	return rootCmd
}

func Execute() {
	cmd := NewRootCmd()

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
