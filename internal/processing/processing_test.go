package processing

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockFileSystem) OpenFile(path string) (io.ReadCloser, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockFileSystem) WriteFile(path string, data []byte) error {
	args := m.Called(path, data)
	return args.Error(0)
}

func (m *MockFileSystem) CreateDir(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

type MockProgressBar struct {
	mock.Mock
}

func (m *MockProgressBar) Increment() {
	m.Called()
}

func (m *MockProgressBar) Finish() {
	m.Called()
}

type ProcessingTestSuite struct {
	suite.Suite
	mockFS          *MockFileSystem
	mockProgressBar *MockProgressBar
	files           []filesystem.FileInfo
	config          FileBatchConfig
}

func (suite *ProcessingTestSuite) SetupTest() {
	suite.mockFS = new(MockFileSystem)
	suite.mockProgressBar = new(MockProgressBar)
	suite.files = []filesystem.FileInfo{
		{Path: "image1.jpg"},
		{Path: "image2.jpg"},
	}
	suite.config = FileBatchConfig{
		Files:       suite.files,
		FS:          suite.mockFS,
		OutputDir:   "output",
		ProgressBar: suite.mockProgressBar,
		Handler: func(task FileTask) error {
			_, err := task.FS.ReadFile(task.File.Path)
			return err
		},
		Concurrency: 1,
	}
}

func (suite *ProcessingTestSuite) TearDownTest() {
	suite.mockFS.AssertExpectations(suite.T())
	suite.mockProgressBar.AssertExpectations(suite.T())
}

func (suite *ProcessingTestSuite) setupSuccessMocks() {
	suite.mockFS.On("ReadFile", "image1.jpg").Return([]byte("content1"), nil)
	suite.mockFS.On("ReadFile", "image2.jpg").Return([]byte("content2"), nil)
	suite.mockProgressBar.On("Increment").Twice()
}

func (suite *ProcessingTestSuite) TestRunFileBatch_ShouldSucceed_WhenAllFilesAreProcessedSuccessfully() {
	suite.setupSuccessMocks()
	errs := RunFileBatch(suite.config)
	suite.Empty(errs)
}

func (suite *ProcessingTestSuite) TestRunFileBatch_ShouldReturnErrors_WhenSomeFilesFail() {
	suite.mockFS.On("ReadFile", "image1.jpg").Return(nil, errors.New("read error"))
	suite.mockFS.On("ReadFile", "image2.jpg").Return([]byte("content2"), nil)
	suite.mockProgressBar.On("Increment").Twice()

	errs := RunFileBatch(suite.config)
	suite.Len(errs, 1)
	suite.Contains(errs[0].Error(), "read error")
	suite.Contains(errs[0].Error(), "image1.jpg")
}

func (suite *ProcessingTestSuite) TestRunFileBatch_ShouldUseDefaultConcurrency_WhenZeroIsProvided() {
	expectedConcurrency := runtime.NumCPU()
	numFiles := expectedConcurrency * 3
	var testFiles []filesystem.FileInfo
	for i := range numFiles {
		testFiles = append(testFiles, filesystem.FileInfo{Path: fmt.Sprintf("img%d.jpg", i)})
	}

	var activeWorkers, maxWorkers int32
	var mu sync.Mutex

	spyProcessor := func(task FileTask) error {
		current := atomic.AddInt32(&activeWorkers, 1)

		mu.Lock()
		if current > maxWorkers {
			maxWorkers = current
		}
		mu.Unlock()

		time.Sleep(10 * time.Millisecond)
		atomic.AddInt32(&activeWorkers, -1)
		return nil
	}

	suite.mockProgressBar.On("Increment").Times(numFiles)

	customConfig := FileBatchConfig{
		Files:       testFiles,
		FS:          suite.mockFS,
		ProgressBar: suite.mockProgressBar,
		Handler:     spyProcessor,
		Concurrency: 0,
	}

	errs := RunFileBatch(customConfig)

	suite.Empty(errs)
	suite.Equal(int32(expectedConcurrency), maxWorkers)
}

func TestProcessingTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessingTestSuite))
}
