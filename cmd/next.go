package cmd

import (
	"errors"
	"fmt"

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
		Run: func(cmd *cobra.Command, _ []string) {
			repository, err := git.PlainOpen(*repoDirectoryPath)
			if err != nil {
				log.Fatalf("failed to open repository: %v", err)
			}
			tagInfo, err := gitutils.GetLatestVersionTag(repository)
			if err != nil {
				log.Warnf("No version tags found: %v", err)
				log.Info("Defaulting to version 1.0.0")
				fmt.Fprintln(cmd.OutOrStdout(), "1.0.0")

				return
			}
			commitMessagesSinceTag, err := gitutils.GetCommitMessagesSinceCommitHash(repository, tagInfo.TagRef.Hash())
			if errors.Is(err, gitutils.ErrNoCommitsFound) {
				log.Info("No commits found since the latest version tag")

				return
			}
			if err != nil {
				log.Fatalf("failed to get commit messages since tag: %v", err)
			}
			nextVersion, err := semverutils.CalculateNextVersion(tagInfo.Name, commitMessagesSinceTag)
			if errors.Is(err, semverutils.ErrNoBump) {
				log.Infof("No bump detected: %v", err)

				return
			}
			if err != nil {
				log.Fatalf("no new version calculated: %v", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), nextVersion)
		},
	}

	return nextCmd
}
