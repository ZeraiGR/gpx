package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	RcFilePerm = os.FileMode(0o644)
	RcDirPerm  = os.FileMode(0o755)
)

type ApplyOptions struct {
	DryRun bool
	Backup bool
}

type ApplyResult struct {
	RCPath      string
	BackupPath  string
	WouldChange bool
	NewContent  string // filled for DryRun (and can be useful for debugging)
}

func ApplyToRC(rcPath string, lines []string, opts ApplyOptions) (*ApplyResult, error) {
	block := RenderBlock(lines)

	old := ""
	if b, err := os.ReadFile(rcPath); err == nil {
		old = string(b)
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("read rc %s: %w", rcPath, err)
	}

	newContent := UpsertBlock(old, block)
	wouldChange := newContent != old

	res := &ApplyResult{
		RCPath:      rcPath,
		WouldChange: wouldChange,
		NewContent:  newContent,
	}

	if opts.DryRun {
		// do not touch filesystem
		return res, nil
	}

	// Ensure dir exists
	dir := filepath.Dir(rcPath)
	if err := os.MkdirAll(dir, RcDirPerm); err != nil {
		return nil, fmt.Errorf("mkdir %s: %w", dir, err)
	}

	// Backup existing file (optional)
	if opts.Backup {
		if _, err := os.Stat(rcPath); err == nil {
			bak := backupPath(rcPath)
			if err := copyFile(rcPath, bak, RcFilePerm); err != nil {
				return nil, fmt.Errorf("backup rc to %s: %w", bak, err)
			}
			res.BackupPath = bak
		} else if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("stat rc %s: %w", rcPath, err)
		}
	}

	// Atomic replace
	tmp := rcPath + ".tmp"
	if err := os.WriteFile(tmp, []byte(newContent), RcFilePerm); err != nil {
		return nil, fmt.Errorf("write temp rc %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, rcPath); err != nil {
		return nil, fmt.Errorf("replace rc %s: %w", rcPath, err)
	}

	return res, nil
}

func backupPath(rcPath string) string {
	// Keep it simple and safe:
	// ~/.zshrc -> ~/.zshrc.gpx.bak (with timestamp to avoid overwriting)
	ts := time.Now().Format("20060102-150405")
	return fmt.Sprintf("%s.gpx.%s.bak", rcPath, ts)
}

func copyFile(src, dst string, perm os.FileMode) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, b, perm)
}
