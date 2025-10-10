// Package cmd provides the command line interface for verscout
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

// NewNextCmd creates and returns a cobra.Command for calculating the next semantic version.
// It uses git operations to find the latest version tag and analyzes commit messages to
// determine the next version according to semantic versioning rules.
func NewNextCmd(git GitInterface, repoDirectoryPath *string) *cobra.Command {
	var noNextVersionExitCode int

	var configPath string

	var firstVersion string

	nextCmd := &cobra.Command{
		Use:   "next",
		Short: "Calculate the next version",
		Long:  "Calculate the next version in the format MAJOR.MINOR.PATCH",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := HandleNextCommand(
				cmd.OutOrStdout(),
				git,
				repoDirectoryPath,
				noNextVersionExitCode,
				configPath,
				firstVersion,
			)
			if err != nil {
				return fmt.Errorf("error while running next command: %w", err)
			}

			return nil
		},
	}

	nextCmd.Flags().
		IntVarP(&noNextVersionExitCode, "exit-code", "e", 0, "The exit code to use when no next version is found")
	nextCmd.Flags().
		StringVarP(&configPath, "config-path", "c", ".verscout-config.yaml", "The path to the verscout config file")
	nextCmd.Flags().
		StringVarP(&firstVersion, "first-version", "f", "1.0.0", "The first version to use if no previous version tags exist")

	return nextCmd
}

// HandleNextCommand performs the version calculation logic for the next command.
// It retrieves the latest version tag, analyzes commit messages since that tag,
// and calculates the next version based on semantic versioning rules.
func HandleNextCommand(
	writer io.Writer,
	git GitInterface,
	repoDirectoryPath *string,
	noNextVersionExitCode int,
	configPath string,
	firstVersion string,
) error {
	config, err := semverutils.LoadBumpConfigFromFile(configPath)
	if err != nil {
		log.Warnf("Failed to load config file: %v", err)
		log.Info("Using default config")

		config = semverutils.DefaultBumpConfig
	}

	log.WithField("configFile", configPath).Info("Using config file")

	repository, err := git.PlainOpen(*repoDirectoryPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	tagInfo, err := gitutils.GetLatestVersionTag(repository)
	if err != nil {
		log.Warnf("No version tags found: %v", err)
		log.WithField("firstVersion", firstVersion).Info("Using provided first version")

		_, err = fmt.Fprintln(writer, firstVersion)
		if err != nil {
			return fmt.Errorf("failed to write version: %w", err)
		}

		return nil
	}

	commitMessagesSinceTag, err := gitutils.GetCommitMessagesSinceCommitHash(repository, tagInfo.Commit.Hash)
	if errors.Is(err, gitutils.ErrNoCommitsFound) {
		log.Infof("No commits found since the latest version tag: %v", err)

		if noNextVersionExitCode != 0 {
			return &ExitError{Code: noNextVersionExitCode, Err: err}
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to get commit messages since tag: %w", err)
	}

	nextVersion, err := semverutils.CalculateNextVersion(tagInfo.Name, commitMessagesSinceTag, config)
	if errors.Is(err, semverutils.ErrNoBump) {
		if noNextVersionExitCode != 0 {
			return &ExitError{Code: noNextVersionExitCode, Err: err}
		}

		log.Infof("No bump detected: %v", err)

		return nil
	}

	if err != nil {
		return fmt.Errorf("no new version calculated: %w", err)
	}

	_, err = fmt.Fprintln(writer, nextVersion)
	if err != nil {
		return fmt.Errorf("failed to write next version: %w", err)
	}

	return nil
}
