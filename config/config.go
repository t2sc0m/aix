package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds aix configuration.
type Config struct {
	Timeout int            `yaml:"timeout"`
	Sandbox string         `yaml:"sandbox"`
	Adapters map[string]AdapterConfig `yaml:"adapters"`
}

// AdapterConfig holds per-adapter settings.
type AdapterConfig struct {
	Enabled bool   `yaml:"enabled"`
	Model   string `yaml:"model"`
}

// Defaults returns hardcoded default configuration.
func Defaults() *Config {
	return &Config{
		Timeout: 300,
		Sandbox: "read-only",
		Adapters: map[string]AdapterConfig{
			"codex": {Enabled: true, Model: ""},
		},
	}
}

// Load reads config from ~/.config/aix/config.yaml.
// Returns defaults if file doesn't exist. Prints warning to stderr if malformed.
func Load() *Config {
	defaults := Defaults()

	home, err := os.UserHomeDir()
	if err != nil {
		return defaults
	}

	configPath := filepath.Join(home, ".config", "aix", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		// File doesn't exist: silent fallback
		return defaults
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to parse %s: %v (using defaults)\n", configPath, err)
		return defaults
	}

	return Merge(defaults, &cfg)
}

// Merge merges file config over defaults. Zero values in file config are ignored.
func Merge(defaults, file *Config) *Config {
	result := *defaults

	if file.Timeout > 0 {
		result.Timeout = file.Timeout
	}
	if file.Sandbox != "" {
		result.Sandbox = file.Sandbox
	}
	if file.Adapters != nil {
		for k, v := range file.Adapters {
			result.Adapters[k] = v
		}
	}

	return &result
}

// Resolve applies CLI flags over config. Empty/zero flag values are ignored.
func Resolve(cfg *Config, flagModel, flagSandbox string, flagTimeout int) (model, sandbox string, timeout int) {
	model = cfg.Adapters["codex"].Model
	if flagModel != "" {
		model = flagModel
	}

	sandbox = cfg.Sandbox
	if flagSandbox != "" {
		sandbox = flagSandbox
	}

	timeout = cfg.Timeout
	if flagTimeout > 0 && flagTimeout != 300 {
		timeout = flagTimeout
	}

	return
}
