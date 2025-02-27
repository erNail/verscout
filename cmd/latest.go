package cmd

import (
	"fmt"
	"io"

	"github.com/erNail/verscout/internal/gitutils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Returns a cobra.Command that can be used as a subcommand.
func NewLatestCmd(git GitInterface, repoDirectoryPath *string) *cobra.Command {
	latestCmd := &cobra.Command{
		Use:   "latest",
		Short: "Scout the latest version tag",
		Long:  "Scout the latest version tag in the format MAJOR.MINOR.PATCH",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := HandleLatestCommand(cmd.OutOrStdout(), git, repoDirectoryPath)
      if err != nil {
        return fmt.Errorf("error while running the latest command: %w", err)
      }
      return nil
		},
	}

	return latestCmd
}

func HandleLatestCommand(writer io.Writer, git GitInterface, repoDirectoryPath *string) error {
  repository, err := git.PlainOpen(*repoDirectoryPath)
  if err != nil {
    return fmt.Errorf("failed to open repository: %w", err)
  }
  semVer, err := gitutils.GetLatestVersion(repository)
  if err != nil {
    log.Warnf("Latest version not found: %v", err)
    return nil
  }
  log.WithField("version", semVer.String()).Info("Found latest version")
  fmt.Fprintln(writer, semVer.String())
  return nil
}