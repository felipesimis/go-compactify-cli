package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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

	mockProcessorFactory := func([]byte) image.ImageProcessor { return nil }
	mockProcessFunc := func(proc image.ImageProcessor) ([]byte, error) {
		return []byte{}, nil
	}

	b.ResetTimer()

	for range b.N {
		err := HandleImageProcessing(ctx, params, stats, mockProcessorFactory, mockProcessFunc)
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

	mockProcessorFactory := func([]byte) image.ImageProcessor { return nil }
	mockProcessFunc := func(proc image.ImageProcessor) ([]byte, error) {
		return []byte{}, nil
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := HandleImageProcessing(ctx, params, stats, mockProcessorFactory, mockProcessFunc)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

type mockHelperFS struct {
	mock.Mock
	filesystem.FileSystem
}

func (m *mockHelperFS) OpenFile(path string) (io.ReadCloser, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *mockHelperFS) WriteFile(path string, data []byte) error {
	args := m.Called(path, data)
	return args.Error(0)
}

type errorReader struct{}

func (e errorReader) Read(p []byte) (n int, err error) { return 0, errors.New("forced read error") }
func (e errorReader) Close() error {
	return nil
}

type ProcessingTestSuite struct {
	suite.Suite
	mockFS *mockHelperFS
}

func (suite *ProcessingTestSuite) SetupTest() {
	suite.mockFS = new(mockHelperFS)
}

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSkip_WhenContextCanceled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	stats := &utils.ImageProcessingStats{}
	params := processing.FileProcessingParams{}

	err := HandleImageProcessing(ctx, params, stats, nil, nil)

	suite.ErrorIs(err, context.Canceled)
	suite.Equal(uint32(1), stats.SkippedImages.Load())
}

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSkip_WhenOpenFileFails() {
	suite.mockFS.On("OpenFile", "test.jpg").Return(nil, errors.New("open error"))

	stats := &utils.ImageProcessingStats{}
	params := processing.FileProcessingParams{
		File: filesystem.FileInfo{Path: "test.jpg", Size: 100},
		FS:   suite.mockFS,
	}

	err := HandleImageProcessing(context.Background(), params, stats, nil, nil)
	suite.ErrorContains(err, "open error")
	suite.Equal(uint32(1), stats.SkippedImages.Load())
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSkip_WhenReadFails() {
	suite.mockFS.On("OpenFile", "corrupt.jpg").Return(errorReader{}, nil)

	stats := &utils.ImageProcessingStats{}
	params := processing.FileProcessingParams{
		File: filesystem.FileInfo{Path: "corrupt.jpg", Size: 100},
		FS:   suite.mockFS,
	}

	err := HandleImageProcessing(context.Background(), params, stats, nil, nil)

	suite.Error(err)
	suite.Contains(err.Error(), "forced read error")
	suite.Equal(uint32(1), stats.SkippedImages.Load())
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSkip_WhenWriteFileFails() {
	suite.mockFS.On("OpenFile", "valid.jpg").Return(io.NopCloser(bytes.NewReader([]byte("data"))), nil)
	suite.mockFS.On("WriteFile", mock.Anything, mock.Anything).Return(errors.New("write error"))

	stats := &utils.ImageProcessingStats{}
	params := processing.FileProcessingParams{
		File:      filesystem.FileInfo{Path: "valid.jpg", Size: 4},
		FS:        suite.mockFS,
		OutputDir: "output",
	}

	mockFactory := func([]byte) image.ImageProcessor { return nil }
	mockProcess := func(proc image.ImageProcessor) ([]byte, error) {
		return []byte("new-data"), nil
	}

	err := HandleImageProcessing(context.Background(), params, stats, mockFactory, mockProcess)
	suite.ErrorContains(err, "write error")
	suite.Equal(uint32(1), stats.SkippedImages.Load())
	suite.mockFS.AssertExpectations(suite.T())
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

func TestProcessingSuite(t *testing.T) {
	suite.Run(t, new(ProcessingTestSuite))
}
