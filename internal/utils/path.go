package utils

import "path/filepath"

func BuildOutputPath(outputDir, relativePath string) string {
	return filepath.Join(outputDir, relativePath)
}
