package docgen

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// NormalizeRule defines a regex pattern and replacement for output normalization.
type NormalizeRule struct {
	Pattern string `yaml:"pattern"`
	Replace string `yaml:"replace"`
}

// Config defines the example generation settings for a documentation section.
type Config struct {
	GameDir   string          `yaml:"game_dir"`
	Normalize []NormalizeRule `yaml:"normalize"`
}

// LoadConfig reads a config.yaml file and returns the parsed Config.
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

// Normalize applies all normalization rules to the input string.
func Normalize(s string, rules []NormalizeRule) string {
	for _, r := range rules {
		re := regexp.MustCompile(r.Pattern)
		s = re.ReplaceAllString(s, r.Replace)
	}
	return strings.TrimSpace(s)
}
