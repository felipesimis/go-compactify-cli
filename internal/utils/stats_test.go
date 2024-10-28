package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImageProcessingStats(t *testing.T) {
	initialSize := uint64(0)
	finalSize := uint64(0)
	skippedImages := uint32(0)
	processedImages := uint32(0)

	stats := NewImageProcessingStats(&initialSize, &finalSize, &skippedImages, &processedImages)

	assert.Equal(t, initialSize, *stats.InitialSize)
	assert.Equal(t, finalSize, *stats.FinalSize)
	assert.Equal(t, skippedImages, *stats.SkippedImages)
	assert.Equal(t, processedImages, *stats.ProcessedImages)
}
