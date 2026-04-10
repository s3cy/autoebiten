package docgen

import (
	"regexp"
	"strings"
)

// NormalizeRule defines a regex pattern and replacement for output normalization.
type NormalizeRule struct {
	Pattern string `yaml:"pattern"`
	Replace string `yaml:"replace"`
}

// Config holds template configuration.
// YAML tags support legacy config.yaml loading during migration.
type Config struct {
	GameDir   string          `yaml:"game_dir"`
	Normalize []NormalizeRule `yaml:"normalize"`
}

// Normalize applies all normalization rules to the input string.
// Note: This will be moved to normalize.go in a future refactor.
func Normalize(s string, rules []NormalizeRule) string {
	for _, r := range rules {
		re := regexp.MustCompile(r.Pattern)
		s = re.ReplaceAllString(s, r.Replace)
	}
	return strings.TrimSpace(s)
}