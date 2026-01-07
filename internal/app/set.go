package app

import (
	"fmt"

	"example.com/gpx/internal/envx"
)

func (a App) SetVars(tokens []string) ([]string, error) {
	vars, err := envx.ParseAssignments(tokens)
	if err != nil {
		return nil, err
	}
	lines, err := vars.ExportLines()
	if err != nil {
		return nil, fmt.Errorf("render exports: %w", err)
	}
	return lines, nil
}
