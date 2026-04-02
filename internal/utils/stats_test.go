package utils

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImageProcessingStats_ConcurrentUpdates(t *testing.T) {
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

	assert.Equal(t, uint32(processNumber), stats.ProcessedImages.Load(), "ProcessedImages should be equal to the number of processes")
	assert.Equal(t, uint64(processNumber*1024), stats.FinalSize.Load(), "FinalSize should be equal to the total size added by all processes")
	assert.Equal(t, uint32(0), stats.SkippedImages.Load(), "SkippedImages should be 0 as we only incremented ProcessedImages")
}
