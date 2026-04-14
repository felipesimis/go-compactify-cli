package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildOutputPath_ShouldReturnCorrectPath(t *testing.T) {
	tests := []struct {
		name string
		outputDir string
		relativePath string
		expected string
	}{
		{
		name: "should return path in root directory",
		outputDir: "output",
		relativePath: "image.jpg",
		expected: "output/image.jpg",
		},
		{
		name: "should return path in subdirectory",
		outputDir: "output",
		relativePath: "subdir/image.jpg",
		expected: "output/subdir/image.jpg",
		},
		{
		name: "should return path with nested subdirectories",
		outputDir: "output",
		relativePath: "subdir/nested/a/b/image.jpg",
		expected: "output/subdir/nested/a/b/image.jpg",
	},
}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildOutputPath(tt.outputDir, tt.relativePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}
