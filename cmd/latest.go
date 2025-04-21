package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/erNail/verscout/internal/gitutils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewLatestCmd creates and returns a cobra.Command for retrieving the latest version tag.
func NewLatestCmd(git GitInterface, repoDirectoryPath *string) *cobra.Command {
	var noLatestVersionExitCode int

	latestCmd := &cobra.Command{
		Use:   "latest",
		Short: "Scout the latest version tag",
		Long:  "Scout the latest version tag in the format MAJOR.MINOR.PATCH",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := HandleLatestCommand(cmd.OutOrStdout(), git, repoDirectoryPath, noLatestVersionExitCode)
			if err != nil {
				return fmt.Errorf("error while running latest command: %w", err)
			}

			return nil
		},
	}

	latestCmd.Flags().
		IntVarP(&noLatestVersionExitCode, "exit-code", "e", 0, "The exit code to use when no latest version is found")

	return latestCmd
}

// HandleLatestCommand performs the version retrieval logic for the latest command.
func HandleLatestCommand(
	writer io.Writer,
	git GitInterface,
	repoDirectoryPath *string,
	noLatestVersionExitCode int,
) error {
	repository, err := git.PlainOpen(*repoDirectoryPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	semVer, err := gitutils.GetLatestVersion(repository)
	if err != nil {
		if errors.Is(err, gitutils.ErrNoTags) || errors.Is(err, gitutils.ErrNoValidVersionTags) {
			if noLatestVersionExitCode != 0 {
				return &ExitError{Code: noLatestVersionExitCode, Err: err}
			}

			log.Warnf("Latest version not found: %v", err)

			return nil
		}

		return fmt.Errorf("failed to get latest version: %w", err)
	}

	log.WithField("version", semVer.String()).Info("Found latest version")

	_, err = fmt.Fprintln(writer, semVer.String())
	if err != nil {
		return fmt.Errorf("failed to write latest version: %w", err)
	}

	return nil
}
