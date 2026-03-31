package processing

import (
	"errors"
	"testing"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/stretchr/testify/assert"
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

func (m *MockFileSystem) WriteFile(path string, data []byte) error {
	args := m.Called(path, data)
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
	params          ProcessFilesParams
}

func (suite *ProcessingTestSuite) SetupTest() {
	suite.mockFS = new(MockFileSystem)
	suite.mockProgressBar = new(MockProgressBar)
	suite.files = []filesystem.FileInfo{
		{Path: "image1.jpg"},
		{Path: "image2.jpg"},
	}
	suite.params = ProcessFilesParams{
		Files:       suite.files,
		FS:          suite.mockFS,
		OutputDir:   "output",
		ProgressBar: suite.mockProgressBar,
		ProcessorFunc: func(params FileProcessingParams) error {
			_, err := params.FS.ReadFile(params.File.Path)
			return err
		},
		Concurrency: 1,
	}
}

func (suite *ProcessingTestSuite) TestProcessFiles() {
	suite.mockFS.On("ReadFile", "image1.jpg").Return([]byte("content1"), nil)
	suite.mockFS.On("ReadFile", "image2.jpg").Return([]byte("content2"), nil)
	suite.mockProgressBar.On("Increment").Twice()

	errs := ProcessFiles(suite.params)
	suite.mockFS.AssertExpectations(suite.T())
	suite.mockProgressBar.AssertExpectations(suite.T())
	suite.Empty(errs)
}

func (suite *ProcessingTestSuite) TestProcessFilesWithError() {
	suite.mockFS.On("ReadFile", "image1.jpg").Return(nil, errors.New("read error"))
	suite.mockFS.On("ReadFile", "image2.jpg").Return([]byte("content2"), nil)
	suite.mockProgressBar.On("Increment").Twice()

	errs := ProcessFiles(suite.params)
	suite.mockFS.AssertExpectations(suite.T())
	suite.mockProgressBar.AssertExpectations(suite.T())

	assert.Len(suite.T(), errs, 1, "Expected one error")
	assert.Contains(suite.T(), errs[0].Error(), "read error", "Expected error message to contain 'read error'")
	assert.Contains(suite.T(), errs[0].Error(), "image1.jpg", "Expected error message to contain 'error processing file image1.jpg'")
}

func TestProcessingTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessingTestSuite))
}
