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
		// Simulates a root directory call
		_ = walkFn("input", filesystem.FileInfo{IsDir: true, RelPath: "."})
		// Simulates a subdirectory
		_ = walkFn("input/folder", filesystem.FileInfo{IsDir: true, RelPath: "folder"})
		// Simulates a valid image file (should be collected)
		_ = walkFn("input/folder/img.jpg", filesystem.FileInfo{IsDir: false, RelPath: "folder/img.jpg"})
		// Simulates a non-image file (should be ignored)
		_ = walkFn("input/folder/doc.txt", filesystem.FileInfo{IsDir: false, RelPath: "folder/doc.txt"})
		return nil
	})

	files, err := DiscoverAndPrepare(suite.mockFS, "input", "output", true)

	suite.NoError(err)
	suite.Len(files, 1, "Should only collect the valid image")
	suite.Equal("folder/img.jpg", files[0].RelPath)
}

func (suite *DiscoveryTestSuite) TestDiscover_Recursive_ShouldSkipOutputDirectory_WhenDestinationIsInsideInput() {
	suite.mockFS.On("Walk", "input", mock.Anything).Return(func(root string, walkFn func(string, filesystem.FileInfo) error) error {
		err := walkFn("output", filesystem.FileInfo{IsDir: true, RelPath: "../output"})
		suite.Equal(filepath.SkipDir, err, "Should skip the output directory to prevent recursive loops")
		return nil
	})

	files, err := DiscoverAndPrepare(suite.mockFS, "input", "output", true)

	suite.NoError(err)
	suite.Empty(files, "Should not collect any files since output directory is skipped")
}

func (suite *DiscoveryTestSuite) TestDiscover_Recursive_ShouldReturnError_WhenFilepathAbsThrows() {
	originalAbs := filepathAbs
	defer func() { filepathAbs = originalAbs }()

	expectedErr := errors.New("abs error")
	filepathAbs = func(path string) (string, error) {
		return "", expectedErr
	}

	files, err := DiscoverAndPrepare(suite.mockFS, "input", "output", true)

	suite.ErrorIs(err, expectedErr)
	suite.Nil(files)
}

func (suite *DiscoveryTestSuite) TestDiscover_Recursive_ShouldReturnError_WhenFilepathAbsFailsInsideWalk() {
	originalAbs := filepathAbs
	defer func() { filepathAbs = originalAbs }()

	expectedErr := errors.New("abs error")
	filepathAbs = func(path string) (string, error) {
		if path == "output" {
			return "/fake/absolute/output", nil
		}
		return "", expectedErr
	}

	suite.mockFS.On("Walk", "input", mock.Anything).Return(func(root string, walkFn func(string, filesystem.FileInfo) error) error {
		return walkFn("input/folder", filesystem.FileInfo{IsDir: true, RelPath: "folder"})
	})

	files, err := DiscoverAndPrepare(suite.mockFS, "input", "output", true)

	suite.ErrorIs(err, expectedErr)
	suite.Nil(files)
}

func (suite *DiscoveryTestSuite) TestDiscover_Recursive_WalkError() {
	expectedErr := errors.New("walk error")
	suite.mockFS.On("Walk", "input", mock.Anything).Return(expectedErr)

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
