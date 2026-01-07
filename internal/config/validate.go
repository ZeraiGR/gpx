package config

import (
	"fmt"
)

func Validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}
	if cfg.Profiles == nil {
		return fmt.Errorf("profiles is missing")
	}
	for pname, vars := range cfg.Profiles {
		if pname == "" {
			return fmt.Errorf("profile name is empty")
		}
		if vars == nil {
			return fmt.Errorf("profile %q has null vars map", pname)
		}
		for k := range vars {
			if err := validateEnvKey(k); err != nil {
				return fmt.Errorf("profile %q: %w", pname, err)
			}
		}
	}
	return nil
}

func validateEnvKey(key string) error {
	if key == "" {
		return fmt.Errorf("empty env key")
	}
	for i, r := range key {
		if i == 0 {
			if !(r == '_' || ('A' <= r && r <= 'Z')) {
				return fmt.Errorf("invalid env key %q: must start with A-Z or _", key)
			}
		} else {
			if !(r == '_' || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9')) {
				return fmt.Errorf("invalid env key %q: only A-Z, 0-9, _ allowed", key)
			}
		}
	}
	return nil
}
