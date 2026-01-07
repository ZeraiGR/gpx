package app

import (
	"fmt"
	"os"

	"github.com/ZeraiGR/gpx/internal/config"
)

type InitResult struct {
	Path   string
	Status string // created / already_exists
}

func InitConfig(path string, force bool) (*InitResult, error) {
	if _, err := os.Stat(path); err == nil && !force {
		return &InitResult{Path: path, Status: "already_exists"}, nil
	}

	cfg := config.DefaultConfig()
	if err := config.Save(path, cfg); err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}
	return &InitResult{Path: path, Status: "created"}, nil
}
