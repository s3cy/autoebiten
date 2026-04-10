package docgen

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig reads a config.yaml file and returns the parsed Config.
// This is a legacy function for the old docgen/docverify commands.
// It will be removed when the new template system is complete.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}