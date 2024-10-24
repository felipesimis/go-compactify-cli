package processing

import (
	"errors"
	"testing"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/stretchr/testify/mock"
)

type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) ReadDir(path string) ([]filesystem.FileInfo, error) {
	args := m.Called(path)
	return args.Get(0).([]filesystem.FileInfo), args.Error(1)
}

func (m *MockFileSystem) CreateSiblingDir(path, suffix string) (string, error) {
	args := m.Called(path, suffix)
	return args.String(0), args.Error(1)
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

func TestProcessFiles(t *testing.T) {
	mockFS := new(MockFileSystem)
	mockProgressBar := new(MockProgressBar)

	files := []filesystem.FileInfo{
		{Path: "image1.jpg"},
		{Path: "image2.jpg"},
	}

	mockFS.On("ReadFile", "image1.jpg").Return([]byte("content1"), nil)
	mockFS.On("ReadFile", "image2.jpg").Return([]byte("content2"), nil)
	mockProgressBar.On("Increment").Twice()

	params := ProcessFilesParams{
		Files:       files,
		FS:          mockFS,
		OutputDir:   "output",
		ProgressBar: mockProgressBar,
		ProcessorFunc: func(params FileProcessingParams) error {
			_, err := params.FS.ReadFile(params.File.Path)
			return err
		},
		Concurrency: 1,
	}
	ProcessFiles(params)

	mockFS.AssertExpectations(t)
	mockProgressBar.AssertExpectations(t)
}

func TestProcessFilesWithError(t *testing.T) {
	mockFS := new(MockFileSystem)
	mockProgressBar := new(MockProgressBar)

	files := []filesystem.FileInfo{
		{Path: "image1.jpg"},
		{Path: "image2.jpg"},
	}

	mockFS.On("ReadFile", "image1.jpg").Return(nil, errors.New("read error"))
	mockFS.On("ReadFile", "image2.jpg").Return([]byte("content2"), nil)
	mockProgressBar.On("Increment").Twice()

	params := ProcessFilesParams{
		Files:       files,
		FS:          mockFS,
		OutputDir:   "output",
		ProgressBar: mockProgressBar,
		ProcessorFunc: func(params FileProcessingParams) error {
			_, err := params.FS.ReadFile(params.File.Path)
			return err
		},
		Concurrency: 1,
	}
	ProcessFiles(params)

	mockFS.AssertExpectations(t)
	mockProgressBar.AssertExpectations(t)
}
