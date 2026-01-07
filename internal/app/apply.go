package app

import (
	"fmt"

	"example.com/gpx/internal/shell"
	"example.com/gpx/internal/state"
)

type ApplyReport struct {
	RCPath      string
	BackupPath  string
	WouldChange bool
	NewContent  string
}

func (a App) ApplyProfileToRC(profile string, rcPath string, opts shell.ApplyOptions) (*ApplyReport, error) {
	lines, err := a.UseProfile(profile)
	if err != nil {
		return nil, err
	}
	res, err := shell.ApplyToRC(rcPath, lines, opts)
	if err != nil {
		return nil, fmt.Errorf("apply to rc: %w", err)
	}

	_ = state.SetActiveProfile(profile)

	return &ApplyReport{
		RCPath:      res.RCPath,
		BackupPath:  res.BackupPath,
		WouldChange: res.WouldChange,
		NewContent:  res.NewContent,
	}, nil
}
