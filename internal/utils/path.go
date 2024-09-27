package utils

import "path/filepath"

func BuildOutputPath(dir string, filename string) string {
	return filepath.Join(dir, filepath.Base(filename))
}
