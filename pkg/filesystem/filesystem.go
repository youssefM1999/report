package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindModuleRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

func CheckValidFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}
	if info.IsDir() {
		return fmt.Errorf("file is a directory: %s", path)
	}
	return nil
}

func ResolvePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	root, err := FindModuleRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, path), nil
}
