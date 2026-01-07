package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Profiles map[string]map[string]string `json:"profiles"`
}

// DefaultPath returns ~/.config/gpx/config.json (на macOS/Linux), windows doesn't supported now
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, ".config", "gpx", "config.json"), nil
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	if err := Validate(&cfg); err != nil {
		return nil, fmt.Errorf("validate config %s: %w", path, err)
	}
	if cfg.Profiles == nil {
		cfg.Profiles = map[string]map[string]string{}
	}
	return &cfg, nil
}

func Save(path string, cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("save config: cfg is nil")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, DirPerm); err != nil {
		return fmt.Errorf("create config dir %s: %w", dir, err)
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, FilePerm); err != nil {
		return fmt.Errorf("write temp config %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("replace config %s: %w", path, err)
	}
	return nil
}

func DefaultConfig() *Config {
	return &Config{
		Profiles: map[string]map[string]string{
			"public": {
				"GOPROXY":     "https://proxy.golang.org,direct",
				"GOPRIVATE":   "",
				"GONOSUMDB":   "",
				"GOTOOLCHAIN": "auto",
			},
		},
	}
}
