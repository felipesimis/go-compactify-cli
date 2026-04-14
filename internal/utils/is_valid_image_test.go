package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidImage_ShouldReturnTrue_WhenExtensionIsSupported(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{name: "jpg", file: "image.jpg"},
		{name: "jpeg", file: "image.jpeg"},
		{name: "png", file: "image.png"},
		{name: "webp", file: "image.webp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, IsValidImage(tt.file))
		})
	}
}

func TestIsValidImage_ShouldReturnFalse_WhenExtensionIsNotSupported(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{name: "gif", file: "image.gif"},
		{name: "bmp", file: "image.bmp"},
		{name: "tiff", file: "image.tiff"},
		{name: "svg", file: "image.svg"},
		{name: "no extension", file: "image"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.False(t, IsValidImage(tt.file))
		})
	}
}
