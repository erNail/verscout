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
			semVer, err := gitutils.GetLatestVersion(repository)
			if err != nil {
				log.Warnf("Latest version not found: %v", err)

				return
			}
			log.WithField("version", semVer.String()).Info("Found latest version")
			fmt.Fprintln(cmd.OutOrStdout(), semVer.String())
		},
	}

	return latestCmd
}
