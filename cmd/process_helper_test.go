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

func (m *mockHelperFS) Walk(root string, walkFn func(path string, info filesystem.FileInfo) error) error {
	args := m.Called(root, walkFn)
	if fn, ok := args.Get(0).(func(string, func(string, filesystem.FileInfo) error) error); ok {
		return fn(root, walkFn)
	}
	return args.Error(0)
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

type fakeProcessor struct {
	resultBytes []byte
	err         error
}

func (f *fakeProcessor) Size() (image.ImageSize, error)         { return image.ImageSize{}, nil }
func (f *fakeProcessor) ImageType() string                      { return "" }
func (f *fakeProcessor) Length() int                            { return 0 }
func (f *fakeProcessor) Metadata() (image.ImageMetadata, error) { return image.ImageMetadata{}, nil }
func (f *fakeProcessor) Process(opts ...image.ProcessOption) ([]byte, error) {
	return f.resultBytes, f.err
}

type ProcessingTestSuite struct {
	suite.Suite
	mockFS *mockHelperFS
}

func (suite *ProcessingTestSuite) SetupTest() {
	suite.mockFS = new(mockHelperFS)
}

func (suite *ProcessingTestSuite) TearDownTest() {
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) defaultProcessingTask() (*utils.ImageProcessingStats, processing.FileTask) {
	stats := &utils.ImageProcessingStats{}
	task := processing.FileTask{
		File: filesystem.FileInfo{Path: "test.jpg", Size: 100},
		FS:   suite.mockFS,
	}
	return stats, task
}

func (suite *ProcessingTestSuite) defaultOperationConfigs() (AppConfig, OperationConfig) {
	appConfig := AppConfig{InputDir: "input"}
	config := OperationConfig{
		FileSystem: suite.mockFS,
		Out:        io.Discard,
	}
	return appConfig, config
}

func (suite *ProcessingTestSuite) TestRunOperation_ShouldWrapFileSystemAndPrintWarning_WhenDryRunIsEnabled() {
	suite.mockFS.On("ReadDir", "input").Return([]filesystem.FileInfo{}, nil)

	out := new(bytes.Buffer)
	appConfig, opCfg := suite.defaultOperationConfigs()
	appConfig.DryRun = true
	opCfg.Out = out

	err := RunOperation(appConfig, opCfg)

	suite.NoError(err)
	suite.Contains(out.String(), "DRY-RUN MODE")
}

func (suite *ProcessingTestSuite) TestRunOperation_ShouldReturnError_WhenReadDirFails() {
	expectedErr := errors.New("directory does not exist")
	suite.mockFS.On("CreateSiblingDir", "input", "").Return("input-suffix", nil).Once()
	suite.mockFS.On("ReadDir", "input").Return(nil, expectedErr).Once()

	appConfig, opCfg := suite.defaultOperationConfigs()
	err := RunOperation(appConfig, opCfg)

	suite.ErrorIs(err, expectedErr)
}

func (suite *ProcessingTestSuite) TestRunOperation_ShouldReturnError_WhenOutputDirResolutionFails() {
	expectedErr := errors.New("permission denied")
	suite.mockFS.On("CreateDir", "/forbidden_output").Return(expectedErr)

	appConfig, opCfg := suite.defaultOperationConfigs()
	appConfig.OutputDir = "/forbidden_output"

	err := RunOperation(appConfig, opCfg)

	suite.ErrorIs(err, expectedErr)
}

func (suite *ProcessingTestSuite) TestRunOperation_ShouldPrintInfoRecursiveMessage_WhenRecursiveFlagIsSet() {
	out := new(bytes.Buffer)
	appConfig, opCfg := suite.defaultOperationConfigs()
	appConfig.Recursive = true
	opCfg.Out = out

	suite.mockFS.On("CreateSiblingDir", "input", "").Return("input-suffix", nil).Once()
	suite.mockFS.On("Walk", "input", mock.Anything).Return(nil).Once()

	err := RunOperation(appConfig, opCfg)

	suite.NoError(err)
	suite.Contains(out.String(), "Analyzing directories recursively in input")
}

func (suite *ProcessingTestSuite) TestRunOperation_ShouldPrintInfoNormalMessage_WhenRecursiveFlagIsDisabled() {
	out := new(bytes.Buffer)
	appConfig, opCfg := suite.defaultOperationConfigs()
	appConfig.Recursive = false
	opCfg.Out = out

	suite.mockFS.On("CreateSiblingDir", "input", "").Return("input-suffix", nil).Once()
	suite.mockFS.On("ReadDir", "input").Return([]filesystem.FileInfo{}, nil).Once()

	err := RunOperation(appConfig, opCfg)

	suite.NoError(err)
	suite.Contains(out.String(), "Analyzing files in input")
}

func (suite *ProcessingTestSuite) TestRunOperation_ShouldSucceed_WhenValidInputs() {
	suite.mockFS.On("ReadDir", "input").Return([]filesystem.FileInfo{{Path: "input/test.jpg", RelPath: "test.jpg", Size: 100, IsDir: false}}, nil)
	suite.mockFS.On("CreateSiblingDir", "input", "-suffix").Return("input-suffix", nil)

	out := new(bytes.Buffer)
	appConfig, opCfg := suite.defaultOperationConfigs()
	opCfg.Out = out
	opCfg.OutputSuffix = "-suffix"
	opCfg.ProcessorFunc = func(ctx context.Context, p processing.FileTask, stats *utils.ImageProcessingStats) error {
		stats.ProcessedImages.Add(1)
		return nil
	}

	err := RunOperation(appConfig, opCfg)

	suite.NoError(err)
	suite.Contains(out.String(), "OPERATION")
	suite.Contains(out.String(), "OUTPUT DIRECTORY")
}

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSkip_WhenContextCanceled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	stats, task := suite.defaultProcessingTask()
	err := HandleImageProcessing(ctx, task, stats, nil, AppConfig{})

	suite.ErrorIs(err, context.Canceled)
	suite.Equal(uint32(1), stats.SkippedImages.Load())
}

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSkip_WhenOpenFileFails() {
	suite.mockFS.On("OpenFile", "test.jpg").Return(nil, errors.New("open error"))

	stats, task := suite.defaultProcessingTask()
	err := HandleImageProcessing(context.Background(), task, stats, nil, AppConfig{})

	suite.ErrorContains(err, "open error")
	suite.Equal(uint32(1), stats.SkippedImages.Load())
}

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSkip_WhenReadFails() {
	suite.mockFS.On("OpenFile", "test.jpg").Return(errorReader{}, nil)

	stats, task := suite.defaultProcessingTask()
	err := HandleImageProcessing(context.Background(), task, stats, nil, AppConfig{})

	suite.Error(err)
	suite.Contains(err.Error(), "forced read error")
	suite.Equal(uint32(1), stats.SkippedImages.Load())
}

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSkip_WhenCreateDirFails() {
	suite.mockFS.On("OpenFile", "test.jpg").Return(io.NopCloser(bytes.NewReader([]byte("data"))), nil)
	suite.mockFS.On("CreateDir", "output").Return(errors.New("create dir error"))

	stats, task := suite.defaultProcessingTask()
	task.OutputDir = "output"

	fakeFactory := func([]byte) image.ImageProcessor {
		return &fakeProcessor{resultBytes: []byte("new-data")}
	}

	err := HandleImageProcessing(context.Background(), task, stats, fakeFactory, AppConfig{})
	suite.ErrorContains(err, "create dir error")
	suite.Equal(uint32(1), stats.SkippedImages.Load())
}

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSkip_WhenWriteFileFails() {
	suite.mockFS.On("OpenFile", "test.jpg").Return(io.NopCloser(bytes.NewReader([]byte("data"))), nil)
	suite.mockFS.On("CreateDir", "output").Return(nil)
	suite.mockFS.On("WriteFile", mock.Anything, mock.Anything).Return(errors.New("write error"))

	stats, task := suite.defaultProcessingTask()
	task.OutputDir = "output"
	fakeFactory := func([]byte) image.ImageProcessor {
		return &fakeProcessor{resultBytes: []byte("new-data"), err: nil}
	}

	err := HandleImageProcessing(context.Background(), task, stats, fakeFactory, AppConfig{})
	suite.ErrorContains(err, "write error")
	suite.Equal(uint32(1), stats.SkippedImages.Load())
}

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSkip_WhenProcessingFails() {
	suite.mockFS.On("OpenFile", "test.jpg").Return(io.NopCloser(bytes.NewReader([]byte("data"))), nil)

	stats, task := suite.defaultProcessingTask()
	fakeFactory := func([]byte) image.ImageProcessor { return &fakeProcessor{err: errors.New("processing error")} }

	err := HandleImageProcessing(context.Background(), task, stats, fakeFactory, AppConfig{})
	suite.ErrorContains(err, "processing error")
	suite.Equal(uint32(1), stats.SkippedImages.Load())
}

func (suite *ProcessingTestSuite) TestHandleImageProcessing_ShouldSucceed_WhenAllStepsPass() {
	originalData, processedData := []byte("original"), []byte("new-data")

	suite.mockFS.On("OpenFile", "test.jpg").Return(io.NopCloser(bytes.NewReader(originalData)), nil)
	suite.mockFS.On("CreateDir", "out").Return(nil)
	suite.mockFS.On("WriteFile", filepath.Join("out", "test.jpg"), processedData).Return(nil)

	stats, task := suite.defaultProcessingTask()
	task.File.Size = int64(len(originalData))
	task.OutputDir = "out"

	fakeFactory := func([]byte) image.ImageProcessor { return &fakeProcessor{resultBytes: processedData} }

	err := HandleImageProcessing(context.Background(), task, stats, fakeFactory, AppConfig{})

	suite.NoError(err)
	suite.Equal(uint32(1), stats.ProcessedImages.Load())
	suite.Equal(uint32(0), stats.SkippedImages.Load())
	suite.Equal(uint64(len(originalData)), stats.InitialSize.Load())
	suite.Equal(uint64(len(processedData)), stats.FinalSize.Load())
}

func (suite *ProcessingTestSuite) TestBuildAppConfigOptions_ShouldReturnCorrectOptions() {
	tests := []struct {
		name           string
		appConfigCfg   AppConfig
		expectedLength int
	}{
		{
			name:           "Should include WithStripMetadata when flag is true",
			appConfigCfg:   AppConfig{StripMetadata: true},
			expectedLength: 1,
		},
		{
			name:           "Should return empty slice when flag is false",
			appConfigCfg:   AppConfig{StripMetadata: false},
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			opts := buildAppOptions(tt.appConfigCfg)
			suite.Equal(tt.expectedLength, len(opts))
		})
	}
}

func (suite *ProcessingTestSuite) TestResolveRootOutputDir_ShouldReturnError_WhenCreateDirFails() {
	suite.mockFS.On("CreateDir", "invalid_output").Return(errors.New("create dir error"))

	appConfig, opCfg := suite.defaultOperationConfigs()
	appConfig.OutputDir = "invalid_output"
	rootOutDir, err := resolveRootOutputDir(appConfig, opCfg)

	suite.ErrorContains(err, "create dir error")
	suite.Empty(rootOutDir)
}

func (suite *ProcessingTestSuite) TestResolveRootOutputDir_ShouldReturnCustomDir_WhenProvidedAndCreated() {
	suite.mockFS.On("CreateDir", "valid_output").Return(nil)

	appConfig, opCfg := suite.defaultOperationConfigs()
	appConfig.OutputDir = "valid_output"

	rootOutDir, err := resolveRootOutputDir(appConfig, opCfg)

	suite.NoError(err)
	suite.Equal("valid_output", rootOutDir)
}

func (suite *ProcessingTestSuite) TestResolveRootOutputDir_ShouldReturnSibling_WhenNoCustomDirProvided() {
	suite.mockFS.On("CreateSiblingDir", "input", "-suffix").Return("input-suffix", nil)

	appConfig, opCfg := suite.defaultOperationConfigs()
	opCfg.OutputSuffix = "-suffix"
	rootOutDir, err := resolveRootOutputDir(appConfig, opCfg)

	suite.NoError(err)
	suite.Equal("input-suffix", rootOutDir)
}

func (suite *ProcessingTestSuite) TestResolveRootOutputDir_ShouldReturnError_WhenCreateSiblingDirFails() {
	suite.mockFS.On("CreateSiblingDir", "input", "-suffix").Return("", errors.New("create sibling dir error"))

	appConfig, opCfg := suite.defaultOperationConfigs()
	opCfg.OutputSuffix = "-suffix"
	rootOutDir, err := resolveRootOutputDir(appConfig, opCfg)

	suite.ErrorContains(err, "create sibling dir error")
	suite.Empty(rootOutDir)
}

func (suite *ProcessingTestSuite) TestResolveImageDestination_ShouldUseModifier_WhenProvided() {
	mockModifier := new(mockPathModifier)
	expectedPath := "custom_output_path.png"
	mockModifier.On("ModifyOutputPath", "original.jpg", "/output").Return(expectedPath)

	task := processing.FileTask{
		File:        filesystem.FileInfo{Path: "/input/original.jpg", RelPath: "original.jpg"},
		OutputDir:   "/output",
		ExtraParams: mockModifier,
	}

	result := resolveImageDestination(task)

	suite.Equal(expectedPath, result)
}

func (suite *ProcessingTestSuite) TestResolveImageDestination_ShouldFallbackToDefault_WhenInterfaceNotImplemented() {
	task := processing.FileTask{
		InputDir:  "/base/input",
		File:      filesystem.FileInfo{Path: "/base/input/subfolder/test.jpg", RelPath: "subfolder/test.jpg"},
		OutputDir: "/base/output",
	}

	result := resolveImageDestination(task)

	expectedPath := filepath.Join("/base/output", "subfolder", "test.jpg")
	suite.Equal(expectedPath, result)
}

func (suite *ProcessingTestSuite) TestResolveImageDestination_ShouldFallbackToBase_WhenRelFails() {
	task := processing.FileTask{
		InputDir:  "/base/input",
		File:      filesystem.FileInfo{Path: "base/../test.jpg", RelPath: ""},
		OutputDir: "/base/output",
	}

	result := resolveImageDestination(task)

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
		err := HandleImageProcessing(ctx, task, stats, mockProcessorFactory, AppConfig{})
		if err != nil {
			b.Fatal(err)
		}
	}
}

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
			err := HandleImageProcessing(ctx, task, stats, mockProcessorFactory, AppConfig{})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
