// Package semverutils provides utilities for handling semantic versioning operations,
// including version parsing, validation, and calculating version bumps based on
// conventional commit messages.
package semverutils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrNoCommitsFound is returned when no git commits are found for analysis.
	ErrNoCommitsFound = errors.New("no commits found")
	// ErrNoBump is returned when no commits are found that would trigger a version bump.
	ErrNoBump = errors.New("no conventional commits found that affect the version")
	// ErrInvalidSemVerTag is returned when a version tag doesn't follow semantic versioning format.
	ErrInvalidSemVerTag = errors.New("invalid semantic version tag")
)

// SemVer represents a semantic version with major, minor, and patch components.
type SemVer struct {
	Major int
	Minor int
	Patch int
}

// BumpType represents the type of version bump to perform.
type BumpType int

// Version bump types ordered by precedence.
const (
	NoBump BumpType = iota
	PatchBump
	MinorBump
	MajorBump
)

// IsValidSemVerTag checks if the provided string is a valid semantic version tag.
// The tag may optionally start with 'v' and must follow the format X.Y.Z where X, Y, and Z are non-negative integers.
func IsValidSemVerTag(semVerString string) bool {
	semVerRegex := regexp.MustCompile(`^v?\d+\.\d+\.\d+$`)

	return semVerRegex.MatchString(semVerString)
}

// ExtractSemVerStruct parses a version tag string and returns a SemVer struct.
// Returns ErrInvalidSemVerTag if the tag does not follow semantic versioning format.
func ExtractSemVerStruct(versionTag string) (*SemVer, error) {
	if !IsValidSemVerTag(versionTag) {
		return nil, ErrInvalidSemVerTag
	}

	version := regexp.MustCompile(`\d+\.\d+\.\d+`).FindString(versionTag)

	parts := strings.Split(version, ".")
	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])

	return &SemVer{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// String returns the string representation of a SemVer in the format X.Y.Z.
func (semVer *SemVer) String() string {
	return fmt.Sprintf("%d.%d.%d", semVer.Major, semVer.Minor, semVer.Patch)
}

// CalculateNextVersion determines the next semantic version based on the current version
// and a list of commit messages following conventional commit format.
// Returns ErrNoCommitsFound if the commit list is empty.
// Returns ErrNoBump if no version bump is required.
// Returns ErrInvalidSemVerTag if the tag does not follow semantic versioning format.
func CalculateNextVersion(versionTag string, commitMessages []string) (string, error) {
	if len(commitMessages) == 0 {
		return "", ErrNoCommitsFound
	}

	semVer, err := ExtractSemVerStruct(versionTag)
	if err != nil {
		// Error type could be ErrInvalidSemVerTag
		return "", fmt.Errorf("failed to extract SemVer struct: %w", err)
	}

	bumpType := determineBumpType(commitMessages)
	if bumpType == NoBump {
		return "", ErrNoBump
	}

	nextSemVer := applyBump(*semVer, bumpType)

	return nextSemVer.String(), nil
}

func determineBumpType(commitMessages []string) BumpType {
	featRegex := regexp.MustCompile(`^feat(\(.*\))?:`)
	fixRegex := regexp.MustCompile(`^fix(\(.*\))?:`)
	breakingChangeRegex := regexp.MustCompile(`(?m)^BREAKING CHANGE:`)

	bumpType := NoBump

	for _, message := range commitMessages {
		switch {
		case breakingChangeRegex.MatchString(message):
			bumpType = MajorBump

			log.Info("Detected bump type: MAJOR")
		case featRegex.MatchString(message):
			if bumpType < MinorBump {
				bumpType = MinorBump

				log.Info("Detected bump type: MINOR")
			}
		case fixRegex.MatchString(message):
			if bumpType < PatchBump {
				bumpType = PatchBump

				log.Info("Detected bump type: PATCH")
			}
		}
	}

	return bumpType
}

func applyBump(semVer SemVer, bumpType BumpType) SemVer {
	switch bumpType {
	case MajorBump:
		semVer.Major++
		semVer.Minor = 0
		semVer.Patch = 0
	case MinorBump:
		semVer.Minor++
		semVer.Patch = 0
	case PatchBump:
		semVer.Patch++
	case NoBump:
		// No changes
	}

	return semVer
}
