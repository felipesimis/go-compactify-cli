package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildOutputPath_ShouldReturnCorrectPath(t *testing.T) {
	tests := []struct {
		name         string
		outputDir    string
		relativePath string
		expected     string
	}{
		{
			name:         "should return path in root directory",
			outputDir:    "output",
			relativePath: "image.jpg",
			expected:     filepath.Join("output", "image.jpg"),
		},
		{
			name:         "should return path in subdirectory",
			outputDir:    "output",
			relativePath: "subdir/image.jpg",
			expected:     filepath.Join("output", "subdir/image.jpg"),
		},
		{
			name:         "should return path with nested subdirectories",
			outputDir:    "output",
			relativePath: "subdir/nested/a/b/image.jpg",
			expected:     filepath.Join("output", "/subdir/nested/a/b/image.jpg"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildOutputPath(tt.outputDir, tt.relativePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetConfigDir_ShouldReturnCorrectConfigDir(t *testing.T) {
	tests := []struct {
		name     string
		home     string
		expected string
	}{
		{
			name:     "should return default config directory in home",
			home:     filepath.FromSlash("/home/user"),
			expected: filepath.Join(filepath.FromSlash("/home/user"), configDirName, appName),
		},
		{
			name:     "should return default config directory in home with different path format",
			home:     filepath.FromSlash("C:\\Users\\user"),
			expected: filepath.Join(filepath.FromSlash("C:\\Users\\user"), configDirName, appName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetConfigDir(tt.home)
			assert.Equal(t, tt.expected, result)
		})
	}
}
