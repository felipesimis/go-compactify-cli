package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidImage(t *testing.T) {
	validFiles := []string{"image.jpg", "image.jpeg", "image.png", "image.webp"}
	invalidFiles := []string{"image.gif", "image.bmp", "image.tiff", "image.svg"}
	for _, file := range validFiles {
		assert.True(t, IsValidImage(file), "Expected %s to be a valid image", file)
	}
	for _, file := range invalidFiles {
		assert.False(t, IsValidImage(file), "Expected %s to be an invalid image", file)
	}
}
