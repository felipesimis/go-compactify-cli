package utils

import "path/filepath"

const (
	configDirName = ".config"
	appName       = "compactify"
)

func BuildOutputPath(outputDir, relativePath string) string {
	return filepath.Join(outputDir, relativePath)
}

func GetConfigDir(home string) string {
	return filepath.Join(home, configDirName, appName)
}
