package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

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

func (m *mockHelperFS) CreateDir(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *mockHelperFS) CreateSiblingDir(inputDir, suffix string) (string, error) {
	args := m.Called(inputDir, suffix)
	return args.Get(0).(string), args.Error(1)
}

func (m *mockHelperFS) ReadDir(path string) ([]filesystem.FileInfo, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]filesystem.FileInfo), args.Error(1)
}

type errorReader struct{}

func (e errorReader) Read(p []byte) (n int, err error) { return 0, errors.New("forced read error") }
func (e errorReader) Close() error {
	return nil
}

type mockPathModifier struct {
	mock.Mock
}

func (m *mockPathModifier) ModifyOutputPath(originalPath, outputDir string) string {
	args := m.Called(originalPath, outputDir)
	return args.String(0)
}

type ProcessingTestSuite struct {
	suite.Suite
	mockFS *mockHelperFS
}

func (suite *ProcessingTestSuite) SetupTest() {
	suite.mockFS = new(mockHelperFS)
}

func (suite *ProcessingTestSuite) TestRunOperation_ShouldWrapFileSystemAndPrintWarning_WhenDryRunIsEnabled() {
	suite.mockFS.On("ReadDir", "input").Return([]filesystem.FileInfo{}, nil)

	out := new(bytes.Buffer)
	global := GlobalConfig{InputDir: "input", DryRun: true}
	opCfg := OperationConfig{FileSystem: suite.mockFS, Out: out}

	err := RunOperation(global, opCfg)

	suite.NoError(err)
	suite.Contains(out.String(), "DRY-RUN MODE: No files will be modified or created on disk.")
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) TestRunOperation_ShouldReturnError_WhenReadDirFails() {
	expectedErr := errors.New("directory does not exist")
	suite.mockFS.On("ReadDir", "input").Return(nil, expectedErr)

	global := GlobalConfig{InputDir: "input"}
	opCfg := OperationConfig{FileSystem: suite.mockFS, Out: io.Discard}

	err := RunOperation(global, opCfg)

	suite.ErrorIs(err, expectedErr)
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) TestRunOperation_ShouldReturnError_WhenOutputDirResolutionFails() {
	suite.mockFS.On("ReadDir", "/input").Return([]filesystem.FileInfo{
		{Path: "/input/test.jpg", Size: 100},
	}, nil)
	expectedErr := errors.New("permission denied")
	suite.mockFS.On("CreateDir", "/forbidden_output").Return(expectedErr)

	out := new(bytes.Buffer)
	global := GlobalConfig{InputDir: "/input", OutputDir: "/forbidden_output"}
	opCfg := OperationConfig{FileSystem: suite.mockFS, Out: out}

	err := RunOperation(global, opCfg)

	suite.ErrorIs(err, expectedErr)
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) TestRunOperation_ShouldSucceed_WhenValidInputs() {
	out := new(bytes.Buffer)
	global := GlobalConfig{InputDir: "input"}

	suite.mockFS.On("ReadDir", "input").Return([]filesystem.FileInfo{{Path: "input/test.jpg", Size: 100}}, nil)
	suite.mockFS.On("CreateSiblingDir", "input", "-suffix").Return("input-suffix", nil)

	opCfg := OperationConfig{
		FileSystem:   suite.mockFS,
		Out:          out,
		OutputSuffix: "-suffix",
		ProcessorFunc: func(ctx context.Context, p processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
			stats.ProcessedImages.Add(1)
			stats.InitialSize.Add(100)
			stats.FinalSize.Add(50)
			return nil
		},
	}

	err := RunOperation(global, opCfg)

	suite.NoError(err)
	suite.Contains(out.String(), "OPERATION")
	suite.Contains(out.String(), "OUTPUT DIRECTORY")
	suite.mockFS.AssertExpectations(suite.T())
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

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSkip_WhenProcessingFails() {
	suite.mockFS.On("OpenFile", "valid.jpg").Return(io.NopCloser(bytes.NewReader([]byte("data"))), nil)

	stats := &utils.ImageProcessingStats{}
	params := processing.FileProcessingParams{
		File: filesystem.FileInfo{Path: "valid.jpg", Size: 4},
		FS:   suite.mockFS,
	}

	mockFactory := func([]byte) image.ImageProcessor { return nil }
	mockProcess := func(proc image.ImageProcessor) ([]byte, error) {
		return nil, errors.New("processing error")
	}

	err := HandleImageProcessing(context.Background(), params, stats, mockFactory, mockProcess)
	suite.ErrorContains(err, "processing error")
	suite.Equal(uint32(1), stats.SkippedImages.Load())
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSucceed_WhenAllStepsPass() {
	originalData := []byte("original")
	processedData := []byte("new-data")

	suite.mockFS.On("OpenFile", "test.jpg").Return(io.NopCloser(bytes.NewReader(originalData)), nil)
	suite.mockFS.On("WriteFile", filepath.Join("out", "test.jpg"), processedData).Return(nil)

	stats := &utils.ImageProcessingStats{}
	params := processing.FileProcessingParams{
		File:      filesystem.FileInfo{Path: "test.jpg", Size: int64(len(originalData))},
		FS:        suite.mockFS,
		OutputDir: "out",
	}

	mockFactory := func([]byte) image.ImageProcessor { return nil }
	mockProcess := func(proc image.ImageProcessor) ([]byte, error) { return processedData, nil }

	err := HandleImageProcessing(context.Background(), params, stats, mockFactory, mockProcess)

	suite.NoError(err)
	suite.Equal(uint32(1), stats.ProcessedImages.Load())
	suite.Equal(uint32(0), stats.SkippedImages.Load())
	suite.Equal(uint64(len(originalData)), stats.InitialSize.Load())
	suite.Equal(uint64(len(processedData)), stats.FinalSize.Load())
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) TestResolveOutputDir_ShouldReturnError_WhenCreateDirFails() {
	suite.mockFS.On("CreateDir", "invalid_output").Return(errors.New("create dir error"))

	global := GlobalConfig{OutputDir: "invalid_output"}
	config := OperationConfig{FileSystem: suite.mockFS}

	out, err := resolveOutputDir(global, config)

	suite.ErrorContains(err, "create dir error")
	suite.Empty(out)
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) TestResolveOutputDir_ShouldReturnCustomDir_WhenProvidedAndCreated() {
	suite.mockFS.On("CreateDir", "valid_output").Return(nil)

	global := GlobalConfig{OutputDir: "valid_output"}
	opCfg := OperationConfig{FileSystem: suite.mockFS}

	out, err := resolveOutputDir(global, opCfg)

	suite.NoError(err)
	suite.Equal("valid_output", out)
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) TestResolveOutputDir_ShouldReturnSibling_WhenNoCustomDirProvided() {
	suite.mockFS.On("CreateSiblingDir", "input", "-suffix").Return("input-suffix", nil)

	global := GlobalConfig{InputDir: "input"}
	opCfg := OperationConfig{FileSystem: suite.mockFS, OutputSuffix: "-suffix"}

	out, err := resolveOutputDir(global, opCfg)

	suite.NoError(err)
	suite.Equal("input-suffix", out)
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) TestResolveOutputDir_ShouldReturnError_WhenCreateSiblingDirFails() {
	suite.mockFS.On("CreateSiblingDir", "input", "-suffix").Return("", errors.New("create sibling dir error"))

	global := GlobalConfig{InputDir: "input"}
	opCfg := OperationConfig{FileSystem: suite.mockFS, OutputSuffix: "-suffix"}

	out, err := resolveOutputDir(global, opCfg)

	suite.ErrorContains(err, "create sibling dir error")
	suite.Empty(out)
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) TestDetermineOutputPath_ShouldUseModifier_WhenProvided() {
	mockModifier := new(mockPathModifier)
	expectedPath := "custom_output_path.png"
	mockModifier.On("ModifyOutputPath", "/input/original.jpg", "/output").Return(expectedPath)

	params := processing.FileProcessingParams{
		File:        filesystem.FileInfo{Path: "/input/original.jpg"},
		OutputDir:   "/output",
		ExtraParams: mockModifier,
	}

	result := determineOutputPath(params)

	suite.Equal(expectedPath, result)
	mockModifier.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) TestDetermineOutputPath_ShouldFallbackToDefault_WhenInterfaceNotImplemented() {
	params := processing.FileProcessingParams{
		InputDir:  "/base/input",
		File:      filesystem.FileInfo{Path: "/base/input/subfolder/test.jpg"},
		OutputDir: "/base/output",
	}

	result := determineOutputPath(params)

	expectedPath := filepath.Join("/base/output", "subfolder", "test.jpg")
	suite.Equal(expectedPath, result)
}

func (suite *ProcessingTestSuite) TestDetermineOutputPath_ShouldFallbackToBase_WhenRelFails() {
	params := processing.FileProcessingParams{
		InputDir:  "/base/input",
		File:      filesystem.FileInfo{Path: "base/../test.jpg"},
		OutputDir: "/base/output",
	}

	result := determineOutputPath(params)

	expectedPath := filepath.Join("/base/output", "test.jpg")
	suite.Equal(expectedPath, result)
}

func (suite *ProcessingTestSuite) TestRenderProcessSummary_ShouldPrintFormattedResults_WhenCalled() {
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
		suite.Run(tt.name, func() {
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
				suite.Contains(printedResult, expectedText)
			}
			for _, notExpectedText := range tt.notExpected {
				suite.NotContains(printedResult, notExpectedText)
			}
		})
	}
}

func TestProcessingSuite(t *testing.T) {
	suite.Run(t, new(ProcessingTestSuite))
}

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
