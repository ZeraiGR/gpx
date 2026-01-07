package app

import (
	"fmt"

	"github.com/ZeraiGR/gpx/internal/envx"
	"github.com/ZeraiGR/gpx/internal/state"
)

func (a App) UseProfile(name string) ([]string, error) {
	cfg, err := a.LoadConfig()
	if err != nil {
		return nil, err
	}
	p, ok := cfg.Profiles[name]
	if !ok {
		return nil, &ProfileNotFoundError{Name: name}
	}

	vars := envx.Vars(p)
	lines, err := vars.ExportLines()
	if err != nil {
		return nil, fmt.Errorf("render exports: %w", err)
	}

	// Mark as active (best-effort; should not break the main command).
	_ = state.SetActiveProfile(name)

	return lines, nil
}
