package cmd

import (
	"context"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
)

func BenchmarkHandleImageProcessingParallel(b *testing.B) {
	ctx := context.Background()

	fs := filesystem.NewFileSystem()
	stats := &utils.ImageProcessingStats{}
	task := processing.FileTask{
		File: filesystem.FileInfo{
			Path: "../test/testdata/large_image_sample.jpg",
			Size: 10 * 1024 * 1024,
		},
		FS:        fs,
		OutputDir: b.TempDir(),
	}

	mockProcessorFactory := func([]byte) image.ImageProcessor { return &fakeProcessor{resultBytes: []byte("processed data")} }

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := HandleImageProcessing(ctx, task, stats, mockProcessorFactory, AppConfig{}, nil)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkHandleImageProcessing(b *testing.B) {
	ctx := context.Background()

	fs := filesystem.NewFileSystem()
	stats := &utils.ImageProcessingStats{}

	task := processing.FileTask{
		File: filesystem.FileInfo{
			Path: "../test/testdata/sample.jpeg",
			Size: 1024,
		},
		FS:        fs,
		OutputDir: b.TempDir(),
	}

	mockProcessorFactory := func([]byte) image.ImageProcessor { return &fakeProcessor{resultBytes: []byte("processed data")} }

	b.ResetTimer()

	for range b.N {
		err := HandleImageProcessing(ctx, task, stats, mockProcessorFactory, AppConfig{}, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}
