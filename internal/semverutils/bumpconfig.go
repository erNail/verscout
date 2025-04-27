package semverutils

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// BumpPatterns holds regex patterns for each bump type.
type BumpPatterns struct {
	MajorPatterns []string `yaml:"majorPatterns"`
	MinorPatterns []string `yaml:"minorPatterns"`
	PatchPatterns []string `yaml:"patchPatterns"`
}

// BumpConfig holds the configuration for version bumping.
type BumpConfig struct {
	Bumps BumpPatterns `yaml:"bumps"`
}

// DefaultBumpConfig provides the verscout default bump patterns.
var DefaultBumpConfig = BumpConfig{
	Bumps: BumpPatterns{
		MajorPatterns: []string{
			`(?m)^BREAKING CHANGE:`,
		},
		MinorPatterns: []string{
			`^feat(\(.*\))?:`,
		},
		PatchPatterns: []string{
			`^fix(\(.*\))?:`,
		},
	},
}

// LoadBumpConfigFromFile loads a BumpConfig from a YAML file.
func LoadBumpConfigFromFile(configFilePath string) (BumpConfig, error) {
	cleanedConfigFilePath := filepath.Clean(configFilePath)

	file, err := os.Open(cleanedConfigFilePath)
	if err != nil {
		return BumpConfig{}, fmt.Errorf("failed to open config file: %w", err)
	}

	var config BumpConfig

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&config)
	if err != nil {
		return BumpConfig{}, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}
