package filesystem

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockFileSystem struct {
	mock.Mock
}

type DryRunFileSystemTestSuite struct {
	suite.Suite
	dryRunFs FileSystem
	mockFS   *MockFileSystem
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

func (suite *DryRunFileSystemTestSuite) SetupTest() {
	suite.mockFS = new(MockFileSystem)
	suite.dryRunFs = NewDryRunFileSystem(suite.mockFS)
}

func (suite *DryRunFileSystemTestSuite) TestReadDir_IsDelegated() {
	suite.mockFS.On("ReadDir", "test/path").Return([]FileInfo{{Path: "image.jpg"}}, nil)
	files, err := suite.dryRunFs.ReadDir("test/path")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), files, 1)
	suite.mockFS.AssertExpectations(suite.T())
}

func TestDryRunFileSystemTestSuite(t *testing.T) {
	suite.Run(t, new(DryRunFileSystemTestSuite))
}
