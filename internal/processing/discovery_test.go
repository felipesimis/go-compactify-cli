package processing

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockDiscoveryFS struct {
	mock.Mock
	filesystem.FileSystem
}

func (m *mockDiscoveryFS) Walk(root string, walkFn func(path string, info filesystem.FileInfo) error) error {
	args := m.Called(root, walkFn)
	if fn, ok := args.Get(0).(func(string, func(string, filesystem.FileInfo) error) error); ok {
		return fn(root, walkFn)
	}
	return args.Error(0)
}

func (m *mockDiscoveryFS) CreateDir(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *mockDiscoveryFS) ReadDir(path string) ([]filesystem.FileInfo, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]filesystem.FileInfo), args.Error(1)
}

type DiscoveryTestSuite struct {
	suite.Suite
	mockFS *mockDiscoveryFS
}

func (suite *DiscoveryTestSuite) SetupTest() {
	suite.mockFS = new(mockDiscoveryFS)
}

func (suite *DiscoveryTestSuite) TearDownTest() {
	suite.mockFS.AssertExpectations(suite.T())
}

func (suite *DiscoveryTestSuite) TestDiscover_Recursive_Success() {
	suite.mockFS.On("Walk", "input", mock.Anything).Return(func(root string, walkFn func(string, filesystem.FileInfo) error) error {
		// Simulates a root directory (should be ignored for creation but not for discovery)
		_ = walkFn("input", filesystem.FileInfo{IsDir: true, RelPath: "."})
		// Simulates a subdirectory (should trigger a CreateDir call)
		_ = walkFn("input/folder", filesystem.FileInfo{IsDir: true, RelPath: "folder"})
		// Simulates a valid image file (should be collected)
		_ = walkFn("input/folder/img.jpg", filesystem.FileInfo{IsDir: false, RelPath: "folder/img.jpg"})
		// Simulates a non-image file (should be ignored)
		_ = walkFn("input/folder/doc.txt", filesystem.FileInfo{IsDir: false, RelPath: "folder/doc.txt"})
		return nil
	})

	suite.mockFS.On("CreateDir", filepath.Join("output", "folder")).Return(nil).Once()

	files, err := DiscoverAndPrepare(suite.mockFS, "input", "output", true)

	suite.NoError(err)
	suite.Len(files, 1, "Should only collect the valid image")
	suite.Equal("folder/img.jpg", files[0].RelPath)
}

func (suite *DiscoveryTestSuite) TestDiscover_Recursive_WalkError() {
	expectedErr := errors.New("walk error")
	suite.mockFS.On("Walk", "input", mock.Anything).Return(expectedErr)

	files, err := DiscoverAndPrepare(suite.mockFS, "input", "output", true)

	suite.ErrorIs(err, expectedErr)
	suite.Nil(files)
}

func (suite *DiscoveryTestSuite) TestDiscover_Recursive_CreateDirError() {
	expectedErr := errors.New("mkdir error")

	suite.mockFS.On("CreateDir", filepath.Join("output", "folder")).Return(expectedErr).Once()
	suite.mockFS.On("Walk", "input", mock.Anything).Return(func(root string, walkFn func(string, filesystem.FileInfo) error) error {
		return walkFn("input/folder", filesystem.FileInfo{IsDir: true, RelPath: "folder"})
	})

	files, err := DiscoverAndPrepare(suite.mockFS, "input", "output", true)

	suite.ErrorIs(err, expectedErr)
	suite.Nil(files)
}

func (suite *DiscoveryTestSuite) TestDiscover_NonRecursive_ReadDirError() {
	expectedErr := errors.New("readdir error")
	suite.mockFS.On("ReadDir", "input").Return(nil, expectedErr).Once()

	files, err := DiscoverAndPrepare(suite.mockFS, "input", "output", false)

	suite.ErrorIs(err, expectedErr)
	suite.Nil(files)
}

func (suite *DiscoveryTestSuite) TestDiscover_NonRecursive_Success() {
	entries := []filesystem.FileInfo{
		{IsDir: true, RelPath: "subdir"},
		{IsDir: false, RelPath: "img1.jpg"},
		{IsDir: false, RelPath: "doc.txt"},
	}

	suite.mockFS.On("ReadDir", "input").Return(entries, nil).Once()

	files, err := DiscoverAndPrepare(suite.mockFS, "input", "output", false)

	suite.NoError(err)
	suite.Len(files, 1, "Should only collect the valid image")
	suite.Equal("img1.jpg", files[0].RelPath)
}

func TestDiscoveryTestSuite(t *testing.T) {
	suite.Run(t, new(DiscoveryTestSuite))
}
