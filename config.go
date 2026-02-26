package docs

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// docsConfig holds user configuration for the docs subcommand.
type docsConfig struct {
	Renderer string `yaml:"renderer"`
}

// configDir returns the directory where the docs config file lives.
// It uses $XDG_CONFIG_HOME/k6 if set, otherwise ~/.config/k6.
func configDir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "k6"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "k6"), nil
}

// loadConfig reads the docs config file and returns the parsed config.
// If the file does not exist, an empty config is returned with no error.
// If the file exists but cannot be parsed, an empty config is returned
// along with the parse error so callers can log a warning.
func loadConfig() (docsConfig, error) {
	dir, err := configDir()
	if err != nil {
		return docsConfig{}, nil
	}

	data, err := os.ReadFile(filepath.Join(dir, "docs.yaml"))
	if err != nil {
		// File not found is perfectly normal — return empty config.
		return docsConfig{}, nil
	}

	var cfg docsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return docsConfig{}, err
	}

	return cfg, nil
}
