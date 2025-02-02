package semverutils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type SemVer struct {
	Major int
	Minor int
	Patch int
}

type BumpType int

const (
	NoBump BumpType = iota
	PatchBump
	MinorBump
	MajorBump
)

func IsValidSemVer(semVerString string) bool {
	semVerRegex := regexp.MustCompile(`^\d+\.\d+\.\d+$`)

	return semVerRegex.MatchString(semVerString)
}

func ExtractSemVerStruct(semVerString string) (*SemVer, error) {
	if !IsValidSemVer(semVerString) {
		return nil, fmt.Errorf("invalid semantic version: %s", semVerString)
	}

	parts := strings.Split(semVerString, ".")
	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])

	return &SemVer{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func (sv *SemVer) String() string {
	return fmt.Sprintf("%d.%d.%d", sv.Major, sv.Minor, sv.Patch)
}

func CalculateNextVersion(currentVersion string, commitMessages []string) (string, error) {
	if len(commitMessages) == 0 {
		return "", errors.New("no commit messages provided")
	}

	semVer, err := ExtractSemVerStruct(currentVersion)
	if err != nil {
		return "", fmt.Errorf("failed to extract SemVer struct: %w", err)
	}

	bumpType := determineBumpType(commitMessages)
	if bumpType == NoBump {
		return "", errors.New("no conventional commits found that affect the version")
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
		case featRegex.MatchString(message):
			if bumpType < MinorBump {
				bumpType = MinorBump
			}
		case fixRegex.MatchString(message):
			if bumpType < PatchBump {
				bumpType = PatchBump
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
