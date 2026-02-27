package docs

import (
	"errors"
	"io/fs"
	"path/filepath"

	"go.k6.io/k6/lib/fsext"
	"gopkg.in/yaml.v3"
)

// docsConfig holds user configuration for the docs subcommand.
type docsConfig struct {
	Renderer string `yaml:"renderer"`
}

// configDir returns the directory where the docs config file lives.
// It uses $XDG_CONFIG_HOME/k6 if set, otherwise ~/.config/k6.
func configDir(env map[string]string) (string, error) {
	if xdg := env["XDG_CONFIG_HOME"]; xdg != "" {
		return filepath.Join(xdg, "k6"), nil
	}
	home := env["HOME"]
	if home == "" {
		return "", errors.New("HOME not set")
	}
	return filepath.Join(home, ".config", "k6"), nil
}

// loadConfig reads the docs config file and returns the parsed config.
// If the file does not exist, an empty config is returned with no error.
// If the file exists but cannot be parsed, an empty config is returned
// along with the parse error so callers can log a warning.
func loadConfig(afs fsext.Fs, env map[string]string) (docsConfig, error) {
	dir, err := configDir(env)
	if err != nil {
		return docsConfig{}, err
	}

	data, err := fsext.ReadFile(afs, filepath.Join(dir, "docs.yaml"))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// File not found is perfectly normal — return empty config.
			return docsConfig{}, nil
		}
		return docsConfig{}, err
	}

	var cfg docsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return docsConfig{}, err
	}

	return cfg, nil
}
