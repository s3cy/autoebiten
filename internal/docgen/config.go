package docgen

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
