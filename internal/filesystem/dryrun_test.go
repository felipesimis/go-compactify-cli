package filesystem

import (
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) ReadDir(path string) ([]FileInfo, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]FileInfo), args.Error(1)
}

func (m *MockFileSystem) CreateDir(path string) error {
	args := m.Called(path)
	return args.Error(0)
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

type DryRunFileSystemTestSuite struct {
	suite.Suite
	dryRunFs FileSystem
	mockFS   *MockFileSystem
}

func (suite *DryRunFileSystemTestSuite) SetupTest() {
	suite.mockFS = new(MockFileSystem)
	suite.dryRunFs = NewDryRunFileSystem(suite.mockFS)
}

func (suite *DryRunFileSystemTestSuite) TestReadDir_ShouldReturnFiles_WhenCalled() {
	suite.mockFS.On("ReadDir", "test/path").Return([]FileInfo{{Path: "image.jpg"}}, nil)
	files, err := suite.dryRunFs.ReadDir("test/path")
	suite.NoError(err)
	suite.Len(files, 1)
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *DryRunFileSystemTestSuite) TestReadFile_ShouldReturnContent_WhenCalled() {
	suite.mockFS.On("ReadFile", "test/file.jpg").Return([]byte("file content"), nil)
	data, err := suite.dryRunFs.ReadFile("test/file.jpg")
	suite.NoError(err)
	suite.Equal([]byte("file content"), data)
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *DryRunFileSystemTestSuite) TestOpenFile_ShouldReturnFile_WhenCalled() {
	suite.mockFS.On("OpenFile", "test/file.jpg").Return(nil, nil)
	_, err := suite.dryRunFs.OpenFile("test/file.jpg")
	suite.NoError(err)
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *DryRunFileSystemTestSuite) TestCreateDir_ShouldReturnNoError_WhenCalled() {
	err := suite.dryRunFs.CreateDir("test/newdir")
	suite.NoError(err)
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *DryRunFileSystemTestSuite) TestWriteFile_ShouldReturnNoError_WhenCalled() {
	err := suite.dryRunFs.WriteFile("test/file.jpg", []byte("file content"))
	suite.NoError(err)
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *DryRunFileSystemTestSuite) TestCreateSiblingDir_ShouldReturnNewPath_WhenCalled() {
	path, err := suite.dryRunFs.CreateSiblingDir("test/input", "-suffix")
	suite.NoError(err)
	suite.Equal("test/input-suffix", filepath.ToSlash(path))
	suite.mockFS.AssertExpectations(suite.T())
}

func TestDryRunFileSystemTestSuite(t *testing.T) {
	suite.Run(t, new(DryRunFileSystemTestSuite))
}
