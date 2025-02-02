package cmd

import (
	"fmt"

	"github.com/erNail/verscout/internal/gitutils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewLatestCmd(git GitInterface, repoDirectoryPath *string) *cobra.Command {
	latestCmd := &cobra.Command{
		Use:   "latest",
		Short: "Scout the latest version tag",
		Long:  "Scout the latest version tag in the format MAJOR.MINOR.PATCH",
		Run: func(cmd *cobra.Command, _ []string) {
			repository, err := git.PlainOpen(*repoDirectoryPath)
			if err != nil {
				log.Fatalf("failed to open repository: %v", err)
			}
			tagInfo, err := gitutils.GetLatestVersionTag(repository)
			if err != nil {
				log.Warn("No version tags found")

				return
			}

			fmt.Fprintln(cmd.OutOrStdout(), tagInfo.Name)
		},
	}

	return latestCmd
}
