package system

import (
	"os"
	"path/filepath"
)

// Setup prepares necessary system-level constructs before the app can run.
func Setup() error {
	cfgDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// create the configuration directory
	err = CreateDir(cfgDir)
	if err != nil {
		return err
	}

	return nil
}

// GetConfigDir returns the directory where application configuration is stored.
func GetConfigDir() (string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "lull"), nil
}

// CreateDir creates a directory and all ancestors leading up to it.
func CreateDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}