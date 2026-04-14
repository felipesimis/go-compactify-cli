package utils

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageProcessingStats_ShouldHandleConcurrentUpdatesCorrectly(t *testing.T) {
	var stats ImageProcessingStats
	var wg sync.WaitGroup

	processNumber := 10000

	for range processNumber {
		wg.Go(func() {
			stats.ProcessedImages.Add(1)
			stats.FinalSize.Add(1024)
		})
	}
	wg.Wait()

	assert.Equal(t, uint32(processNumber), stats.ProcessedImages.Load())
	assert.Equal(t, uint64(processNumber*1024), stats.FinalSize.Load())
	assert.Equal(t, uint32(0), stats.SkippedImages.Load())
}

func TestImageProcessingStats_ShouldTrackInitialSize(t *testing.T) {
	var stats ImageProcessingStats
	initialValue := uint64(5000)

	stats.InitialSize.Store(initialValue)
	assert.Equal(t, initialValue, stats.InitialSize.Load())
}
