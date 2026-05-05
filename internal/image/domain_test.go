package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGravity_IsValid_ShouldEnforceDomainBounds(t *testing.T) {
	tests := []struct {
		name     string
		input    Gravity
		expected bool
	}{
		{"Valid lower bound", GravityCentre, true},
		{"Valid upper bound", GravitySmart, true},
		{"Valid middle value", GravityEast, true},
		{"Invalid negative value", Gravity(-1), false},
		{"Invalid above max bound", maxGravity, false},
		{"Invalid arbitrary high value", Gravity(99), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProcessOptions_ShouldMutateDomainOptionsCorrectly(t *testing.T) {
	tests := []struct {
		name     string
		option   ProcessOption
		expected domainOptions
	}{
		{
			name:     "WithResize should set width and height",
			option:   WithResize(800, 600),
			expected: domainOptions{width: 800, height: 600},
		},
		{
			name:     "WithConvert should set format",
			option:   WithConvert("webp"),
			expected: domainOptions{format: "webp"},
		},
		{
			name:     "WithCrop should set crop flag, dimensions and gravity",
			option:   WithCrop(400, 400, GravityCentre),
			expected: domainOptions{crop: true, width: 400, height: 400, gravity: GravityCentre},
		},
		{
			name:     "WithGrayscale should set grayscale flag",
			option:   WithGrayscale(),
			expected: domainOptions{grayscale: true},
		},
		{
			name:     "WithFlip should set flip flag",
			option:   WithFlip(),
			expected: domainOptions{flip: true},
		},
		{
			name:     "WithEnlarge should set enlarge flag and dimensions",
			option:   WithEnlarge(2000, 2000),
			expected: domainOptions{enlarge: true, width: 2000, height: 2000},
		},
		{
			name:     "WithThumbnail should set thumbnail width",
			option:   WithThumbnail(250),
			expected: domainOptions{thumbnailWidth: 250},
		},
		{
			name:     "WithPalette should set palette flag",
			option:   WithPalette(),
			expected: domainOptions{palette: true},
		},
		{
			name:     "WithLosslessCompress should set lossless flag",
			option:   WithLosslessCompress(),
			expected: domainOptions{lossless: true},
		},
		{
			name:     "WithStripMetadata should set stripMetadata flag",
			option:   WithStripMetadata(),
			expected: domainOptions{stripMetadata: true},
		},
		{
			name:     "WithQuality should set quality value",
			option:   WithQuality(85),
			expected: domainOptions{quality: 85},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &domainOptions{}
			tt.option(opts)
			assert.Equal(t, tt.expected, *opts)
		})
	}
}

func TestProcessOptions_ShouldChainSuccessfully(t *testing.T) {
	opts := &domainOptions{}

	WithResize(1024, 768)(opts)
	WithConvert("png")(opts)
	WithGrayscale()(opts)

	assert.Equal(t, 1024, opts.width)
	assert.Equal(t, 768, opts.height)
	assert.Equal(t, "png", opts.format)
	assert.True(t, opts.grayscale)
}
