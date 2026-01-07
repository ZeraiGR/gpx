package app

import (
	"fmt"

	"github.com/ZeraiGR/gpx/internal/config"
	"github.com/ZeraiGR/gpx/internal/envx"
)

func (a App) AddProfile(name string) error {
	cfg, err := a.LoadConfig()
	if err != nil {
		return err
	}
	if name == "" {
		return fmt.Errorf("profile name is empty")
	}
	if _, exists := cfg.Profiles[name]; exists {
		return fmt.Errorf("profile %q already exists", name)
	}
	cfg.Profiles[name] = map[string]string{}
	if err := config.Save(a.ConfigPath, cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}

func (a App) RemoveProfile(name string) error {
	cfg, err := a.LoadConfig()
	if err != nil {
		return err
	}
	if _, exists := cfg.Profiles[name]; !exists {
		return &ProfileNotFoundError{Name: name}
	}
	delete(cfg.Profiles, name)
	if err := config.Save(a.ConfigPath, cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}

func (a App) RenameProfile(oldName, newName string) error {
	cfg, err := a.LoadConfig()
	if err != nil {
		return err
	}
	if _, ok := cfg.Profiles[oldName]; !ok {
		return &ProfileNotFoundError{Name: oldName}
	}
	if newName == "" {
		return fmt.Errorf("new profile name is empty")
	}
	if _, exists := cfg.Profiles[newName]; exists {
		return fmt.Errorf("profile %q already exists", newName)
	}
	cfg.Profiles[newName] = cfg.Profiles[oldName]
	delete(cfg.Profiles, oldName)

	if err := config.Save(a.ConfigPath, cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}

func (a App) ShowProfile(name string) (map[string]string, error) {
	cfg, err := a.LoadConfig()
	if err != nil {
		return nil, err
	}
	p, ok := cfg.Profiles[name]
	if !ok {
		return nil, &ProfileNotFoundError{Name: name}
	}
	return p, nil
}

func (a App) SetProfileVars(profile string, tokens []string) error {
	cfg, err := a.LoadConfig()
	if err != nil {
		return err
	}
	p, ok := cfg.Profiles[profile]
	if !ok {
		return &ProfileNotFoundError{Name: profile}
	}

	vars, err := envx.ParseAssignments(tokens)
	if err != nil {
		return err
	}

	for k, v := range vars {
		p[k] = v
	}

	if err := config.Save(a.ConfigPath, cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}

func (a App) UnsetProfileVars(profile string, keys []string) error {
	cfg, err := a.LoadConfig()
	if err != nil {
		return err
	}
	p, ok := cfg.Profiles[profile]
	if !ok {
		return &ProfileNotFoundError{Name: profile}
	}
	for _, k := range keys {
		k = envx.NormalizeKey(k)
		if err := envx.ValidateKey(k); err != nil {
			return err
		}
		delete(p, k)
	}
	if err := config.Save(a.ConfigPath, cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}
