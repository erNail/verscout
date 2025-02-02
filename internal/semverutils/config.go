package semverutils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type BumpConfig struct {
	MajorPatterns []string `yaml:"majorPatterns"`
	MinorPatterns []string `yaml:"minorPatterns"`
	PatchPatterns []string `yaml:"patchPatterns"`
}

var DefaultConfig = BumpConfig{
	MajorPatterns: []string{`(?m)^BREAKING CHANGE:`},
	MinorPatterns: []string{`^feat(\(.*\))?:`},
	PatchPatterns: []string{`^fix(\(.*\))?:`},
}

func LoadConfigFromFile(filePath string) (BumpConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return BumpConfig{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config BumpConfig

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return BumpConfig{}, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}
