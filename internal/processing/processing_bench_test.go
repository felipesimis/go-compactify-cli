package processing

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
)

type benchProgressBar struct{}

func (b benchProgressBar) Increment() {}
func (b benchProgressBar) Finish()    {}

func BenchmarkProcessFiles_ConcurrencyScaling(b *testing.B) {
	var files []filesystem.FileInfo
	for i := range 10000 {
		files = append(files, filesystem.FileInfo{
			Path:    fmt.Sprintf("/fake/input/img_%d.jpg", i),
			RelPath: fmt.Sprintf("img_%d.jpg", i),
			Size:    2048 * 1024,
		})
	}

	concurrencyLevels := []int{1, 2, 4, 8, 16, 32, 64}

	for _, workers := range concurrencyLevels {
		b.Run(fmt.Sprintf("Workers-%d", workers), func(b *testing.B) {

			var processedCount int32
			mockProcessor := func(params FileProcessingParams) error {
				time.Sleep(1 * time.Millisecond)
				atomic.AddInt32(&processedCount, 1)
				return nil
			}

			params := ProcessFilesParams{
				Files:         files,
				ProgressBar:   benchProgressBar{},
				ProcessorFunc: mockProcessor,
				Concurrency:   workers,
			}

			for b.Loop() {
				atomic.StoreInt32(&processedCount, 0)

				errs := ProcessFiles(params)

				if len(errs) > 0 {
					b.Fatalf("Unexpected errors: %v", errs)
				}
				if atomic.LoadInt32(&processedCount) != 10000 {
					b.Fatalf("Expected 10,000 files processed, got %d", processedCount)
				}
			}
		})
	}
}
