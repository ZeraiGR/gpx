package app

import (
	"fmt"

	"github.com/ZeraiGR/gpx/internal/config"
)

type App struct {
	ConfigPath string
}

func (a App) LoadConfig() (*config.Config, error) {
	cfg, err := config.Load(a.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return cfg, nil
}
