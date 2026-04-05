package testkit

import (
	"os"
	"time"
)

// config holds the configuration for Launch.
type config struct {
	timeout time.Duration
	args    []string
	env     map[string]string
}

// Option configures Launch behavior.
type Option func(*config)

// WithTimeout sets the timeout for game operations.
// Default is 30 seconds.
func WithTimeout(d time.Duration) Option {
	return func(c *config) {
		c.timeout = d
	}
}

// WithArgs sets additional command-line arguments for the game binary.
func WithArgs(args ...string) Option {
	return func(c *config) {
		c.args = append(c.args, args...)
	}
}

// WithEnv sets environment variables for the game process.
// These are added to the current environment, not replacing it.
func WithEnv(key, value string) Option {
	return func(c *config) {
		if c.env == nil {
			c.env = make(map[string]string)
		}
		c.env[key] = value
	}
}

// defaultConfig returns a config with sensible defaults.
func defaultConfig() *config {
	return &config{
		timeout: 30 * time.Second,
		args:    []string{},
		env:     make(map[string]string),
	}
}

// buildEnv constructs the environment slice for the game process.
// It merges the current environment with the configured extras.
func buildEnv(c *config, socketPath string) []string {
	// Start with current environment
	env := os.Environ()

	// Add configured environment variables
	for key, value := range c.env {
		env = append(env, key+"="+value)
	}

	// Add AUTOEBITEN_SOCKET for unique socket path
	env = append(env, "AUTOEBITEN_SOCKET="+socketPath)

	return env
}
