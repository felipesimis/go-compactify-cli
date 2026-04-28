package cmd

import (
	"context"
	"fmt"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/stretchr/testify/assert"
)

func BenchmarkHandleImageProcessing(b *testing.B) {
	ctx := context.Background()

	fs := filesystem.NewFileSystem()
	stats := &utils.ImageProcessingStats{}

	params := processing.FileProcessingParams{
		File: filesystem.FileInfo{
			Path: "../test/testdata/sample.jpeg",
			Size: 1024,
		},
		FS:        fs,
		OutputDir: b.TempDir(),
	}

	mockProcessFunc := func(proc image.ImageProcessor) ([]byte, error) {
		return []byte{}, nil
	}

	b.ResetTimer()

	for range b.N {
		err := HandleImageProcessing(ctx, params, stats, mockProcessFunc)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHandleImageProcessingParallel(b *testing.B) {
	ctx := context.Background()

	fs := filesystem.NewFileSystem()
	stats := &utils.ImageProcessingStats{}
	params := processing.FileProcessingParams{
		File: filesystem.FileInfo{
			Path: "../test/testdata/large_image_sample.jpg",
			Size: 10 * 1024 * 1024,
		},
		FS:        fs,
		OutputDir: b.TempDir(),
	}

	mockProcessFunc := func(proc image.ImageProcessor) ([]byte, error) {
		return []byte{}, nil
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := HandleImageProcessing(ctx, params, stats, mockProcessFunc)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestRenderProcessSummary_ShouldPrintFormattedResults_WhenCalled(t *testing.T) {
	tests := []struct {
		name            string
		skippedImages   uint32
		processedImages uint32
		errors          []error
		expected        []string
		notExpected     []string
	}{
		{
			name:            "ShouldHideSkippedRow_WhenSkippedImagesIsZero",
			skippedImages:   0,
			processedImages: 10,
			errors:          nil,
			expected: []string{
				"OPERATION", "IMPACT", "OUTPUT DIRECTORY",
				"10 images", "0", "10",
				"10.00 MB", "5.00 MB", "50.00%",
				"output",
			},
			notExpected: []string{
				"ERRORS DETECTED",
				"Skipped",
			},
		},
		{
			name:            "ShouldRenderErrorInfo_WhenErrorsArePresent",
			skippedImages:   3,
			processedImages: 7,
			errors: []error{
				fmt.Errorf("file 'fake.jpg': read error"),
				fmt.Errorf("permission denied"),
			},
			expected: []string{
				"2 ERRORS DETECTED",
				"fake.jpg",
				"read error",
				"permission denied",
			},
		},
		{
			name:            "ShouldShowSkippedRow_WhenSkippedImagesIsGreaterThanZero",
			skippedImages:   2,
			processedImages: 5,
			errors:          nil,
			expected: []string{
				"OPERATION",
				"Skipped",
				"2",
				"5",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := utils.NewResultBuilder(utils.RealTimeProvider{})
			rb.
				SetOriginalBytes(10485760). // 10 MB
				SetProcessedBytes(5242880). // 5 MB
				SetTotalImages(10).
				SetSkippedImages(tt.skippedImages).
				SetProcessedImages(tt.processedImages).
				SetOutputDirectory("output")

			if tt.errors != nil {
				rb.SetErrors(tt.errors)
			}
			result := rb.Build()
			printedResult := RenderProcessSummary(result)

			for _, expectedText := range tt.expected {
				assert.Contains(t, printedResult, expectedText)
			}
			for _, notExpectedText := range tt.notExpected {
				assert.NotContains(t, printedResult, notExpectedText)
			}
		})
	}
}
