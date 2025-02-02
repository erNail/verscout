package cmd

import (
	"fmt"

	"github.com/erNail/verscout/internal/gitutils"
	"github.com/erNail/verscout/internal/semverutils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
				log.Warn("No version tags found")
				fmt.Fprintln(cmd.OutOrStdout(), "1.0.0")

				return
			}
			commitMessagesSinceTag, err := gitutils.GetCommitMessagesSinceCommitHash(repository, tagInfo.TagRef.Hash())
			if err != nil {
				log.Fatalf("failed to get commit messages since tag: %v", err)
			}
			nextVersion, err := semverutils.CalculateNextVersion(tagInfo.Name, commitMessagesSinceTag)
			if err != nil {
				log.Warnf("no new version calculated: %v", err)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), nextVersion)
			}
		},
	}

	return nextCmd
}
