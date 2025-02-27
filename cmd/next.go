package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/erNail/verscout/internal/gitutils"
	"github.com/erNail/verscout/internal/semverutils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Returns a cobra.Command that can be used as a subcommand.
func NewNextCmd(git GitInterface, repoDirectoryPath *string) *cobra.Command {
	nextCmd := &cobra.Command{
		Use:   "next",
		Short: "Calculate the next version",
		Long:  "Calculate the next version in the format MAJOR.MINOR.PATCH",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := HandleNextCommand(cmd.OutOrStdout(), git, repoDirectoryPath)
			if err != nil {
				return fmt.Errorf("error while running next command: %w", err)
			}
			return nil
		},
	}

	return nextCmd
}

func HandleNextCommand(writer io.Writer, git GitInterface, repoDirectoryPath *string) error {
	repository, err := git.PlainOpen(*repoDirectoryPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}
	tagInfo, err := gitutils.GetLatestVersionTag(repository)
	if err != nil {
		defaultVersion := "1.0.0"
		log.Warnf("No version tags found: %v", err)
		log.WithField("defaultVersion", defaultVersion).Info("Using default version")
		fmt.Fprintln(writer, defaultVersion)
		return nil
	}
	commitMessagesSinceTag, err := gitutils.GetCommitMessagesSinceCommitHash(repository, tagInfo.TagRef.Hash())
	if errors.Is(err, gitutils.ErrNoCommitsFound) {
		log.Infof("No commits found since the latest version tag: %v", err)

		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get commit messages since tag: %w", err)
	}
	nextVersion, err := semverutils.CalculateNextVersion(tagInfo.Name, commitMessagesSinceTag)
	if errors.Is(err, semverutils.ErrNoBump) {
		log.Infof("No bump detected: %v", err)

		return nil
	}
	if err != nil {
		return fmt.Errorf("no new version calculated: %w", err)
	}
	fmt.Fprintln(writer, nextVersion)
	return nil
}