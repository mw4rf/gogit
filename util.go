package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Returns the user configuration directory
// On Unix systems, it is usually ~/.config/gogit
// On Windows, it is usually %APPDATA%\gogit
// On macOS, it is usually ~/Library/Application Support/gogit
// The directory is created if it does not exist
func GetUserConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to HOME environment variable
		configDir = os.Getenv("HOME")
		if configDir == "" {
			// Fallback to the current working directory if HOME is not set
			configDir, _ = os.Getwd()
		}
		configDir = filepath.Join(configDir, ".config")
	}
	configDir = filepath.Join(configDir, "gogit")

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		// Log the error and proceed with the directory path
		fmt.Println(ColorOutput(ColorRed, fmt.Sprintf("Could not create directory %s: %s\n", configDir, err)))
		os.Exit(1)
	}

	return configDir
}

