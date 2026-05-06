package filesystem

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (suite *FileSystemTestSuite) SetupTest() {
	suite.mockOS = new(MockOSOperations)
	suite.mockDir = new(MockDir)
	suite.mockFile = new(MockFile)
	suite.fs = &FileSystemWrapper{os: suite.mockOS}
	suite.path = "/mock/dir"
}

func (suite *FileSystemTestSuite) TearDownTest() {
	suite.mockOS.AssertExpectations(suite.T())
	suite.mockDir.AssertExpectations(suite.T())
	suite.mockFile.AssertExpectations(suite.T())
}

func (suite *FileSystemTestSuite) TestNewFileSystem_ShouldReturnFileSystemWrapper() {
	fs := NewFileSystem()
	suite.NotNil(fs)
	_, ok := fs.(*FileSystemWrapper)
	suite.True(ok)
}

func (suite *FileSystemTestSuite) TestReadDir_ShouldReturnFiles_WhenDirectoryIsValid() {
	files := []os.FileInfo{
		FakeFileInfo{name: "image1.jpg", size: 1024, isDir: false},
		FakeFileInfo{name: "image2.jpeg", size: 2048, isDir: false},
		FakeFileInfo{name: "image3.png", size: 4096, isDir: false},
		FakeFileInfo{name: "image4.webp", size: 8192, isDir: false},
		FakeFileInfo{name: "file1.txt", size: 2048, isDir: false},
		FakeFileInfo{name: "subdir", size: 0, isDir: true},
	}

	suite.mockOS.On("Open", suite.path).Return(suite.mockFile, nil)
	suite.mockFile.On("Readdir", -1).Return(files, nil)
	suite.mockFile.On("Close").Return(nil)

	result, err := suite.fs.ReadDir(suite.path)
	suite.NoError(err)
	suite.Len(result, 6)
}

func (suite *FileSystemTestSuite) TestReadDir_ShouldReturnErrOpenDir_WhenOpenFails() {
	suite.mockOS.On("Open", suite.path).Return(nil, errors.New("simulated open error"))

	result, err := suite.fs.ReadDir(suite.path)
	expectedErr := &ErrOpenDir{Path: suite.path, Err: errors.New("simulated open error")}
	suite.Nil(result)
	suite.EqualError(err, expectedErr.Error())
}

func (suite *FileSystemTestSuite) TestReadDir_ShouldReturnErrReadDir_WhenReaddirFails() {
	suite.mockOS.On("Open", suite.path).Return(suite.mockFile, nil)
	suite.mockFile.On("Readdir", -1).Return(nil, errors.New("simulated readdir error"))
	suite.mockFile.On("Close").Return(nil)

	result, err := suite.fs.ReadDir(suite.path)
	expectedErr := &ErrReadDir{Path: suite.path, Err: errors.New("simulated readdir error")}
	suite.Nil(result)
	suite.EqualError(err, expectedErr.Error())
}

func (suite *FileSystemTestSuite) TestCreateDir_ShouldReturnErrCreateDir_WhenMkdirAllFails() {
	suite.mockOS.On("MkdirAll", suite.path, os.ModePerm).Return(errors.New("mock error"))

	err := suite.fs.CreateDir(suite.path)
	expectedErr := &ErrCreateDir{Path: suite.path, Err: errors.New("mock error")}
	suite.EqualError(err, expectedErr.Error())
}

func (suite *FileSystemTestSuite) TestCreateDir_ShouldReturnNoError_WhenMkdirAllSucceeds() {
	suite.mockOS.On("MkdirAll", suite.path, os.ModePerm).Return(nil)

	err := suite.fs.CreateDir(suite.path)
	suite.NoError(err)
}

func (suite *FileSystemTestSuite) TestCreateSiblingDir_ShouldReturnErrCreateSiblingDir_WhenMkdirFails() {
	expectedPath := suite.path + "-suffix"
	suite.mockOS.On("Mkdir", expectedPath, os.ModePerm).Return(errors.New("mock error"))

	newDir, err := suite.fs.CreateSiblingDir(suite.path, "-suffix")
	expectedErr := &ErrCreateSiblingDir{Path: suite.path, Err: errors.New("mock error")}
	suite.Empty(newDir)
	suite.EqualError(err, expectedErr.Error())
}

func (suite *FileSystemTestSuite) TestCreateSiblingDir_ShouldReturnNewPath_WhenMkdirSucceeds() {
	expectedPath := suite.path + "-suffix"
	suite.mockOS.On("Mkdir", expectedPath, os.ModePerm).Return(nil)

	newDir, err := suite.fs.CreateSiblingDir(suite.path, "-suffix")
	suite.NoError(err)
	suite.Equal(expectedPath, newDir)
}

func (suite *FileSystemTestSuite) TestReadFile_ShouldReturnErrReadFile_WhenReadFails() {
	expectedErr := &ErrReadFile{Path: suite.path, Err: errors.New("mock error")}
	suite.mockOS.On("ReadFile", suite.path).Return(nil, expectedErr.Err)

	data, err := suite.fs.ReadFile(suite.path)
	suite.Nil(data)
	suite.EqualError(err, expectedErr.Error())
}

func (suite *FileSystemTestSuite) TestReadFile_ShouldReturnContent_WhenReadSucceeds() {
	expectedData := []byte("file content")
	suite.mockOS.On("ReadFile", suite.path).Return(expectedData, nil)

	data, err := suite.fs.ReadFile(suite.path)
	suite.NoError(err)
	suite.Equal(expectedData, data)
}

func (suite *FileSystemTestSuite) TestOpenFile_ShouldReturnErrReadFile_WhenOpenFails() {
	expectedErr := &ErrReadFile{Path: suite.path, Err: errors.New("mock error")}
	suite.mockOS.On("Open", suite.path).Return(nil, expectedErr.Err)

	file, err := suite.fs.OpenFile(suite.path)
	suite.Nil(file)
	suite.EqualError(err, expectedErr.Error())
}

func (suite *FileSystemTestSuite) TestOpenFile_ShouldReturnFile_WhenOpenSucceeds() {
	suite.mockOS.On("Open", suite.path).Return(suite.mockFile, nil)

	file, err := suite.fs.OpenFile(suite.path)
	suite.NoError(err)
	suite.Equal(suite.mockFile, file)
}

func (suite *FileSystemTestSuite) TestWriteFile_ShouldReturnErrWriteFile_WhenWriteFails() {
	expectedErr := &ErrWriteFile{Path: suite.path, Err: errors.New("mock error")}
	suite.mockOS.On("WriteFile", suite.path, []byte("data"), os.FileMode(0644)).Return(expectedErr.Err)

	err := suite.fs.WriteFile(suite.path, []byte("data"))
	suite.EqualError(err, expectedErr.Error())
}

func (suite *FileSystemTestSuite) TestWriteFile_ShouldReturnNoError_WhenWriteSucceeds() {
	data := []byte("data")
	suite.mockOS.On("WriteFile", suite.path, data, os.FileMode(0644)).Return(nil)

	err := suite.fs.WriteFile(suite.path, data)
	suite.NoError(err)
}

func (suite *FileSystemTestSuite) TestWalk_CallbackError() {
	expectedErr := errors.New("callback error")

	suite.mockOS.On("Walk", suite.path, mock.Anything).Return(expectedErr).Run(func(args mock.Arguments) {
		walkFn := args.Get(1).(func(string, os.DirEntry, error) error)
		_ = walkFn(suite.path+"/file.jpg", nil, expectedErr)
	}).Once()

	err := suite.fs.Walk(suite.path, func(path string, info FileInfo) error {
		return nil
	})
	suite.Error(err)
	suite.ErrorIs(err, expectedErr)
	suite.EqualError(err, (&ErrWalk{Path: suite.path, Err: expectedErr}).Error())
}

func (suite *FileSystemTestSuite) TestWalk_ShouldProcessAllEntries() {
	tests := []struct {
		name     string
		isDir    bool
		fileName string
	}{
		{"ProcessesDirectory", true, "subdir"},
		{"ProcessesNonImageFiles", false, "notes.txt"},
		{"ProcessesImages", false, "photo.jpg"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			fullPath := suite.path + "/" + tt.fileName

			suite.mockOS.On("Walk", suite.path, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				walkFn := args.Get(1).(func(string, os.DirEntry, error) error)

				suite.mockDir.On("IsDir").Return(tt.isDir).Once()
				suite.mockDir.On("Info").Return(FakeFileInfo{size: 100, isDir: tt.isDir}, nil).Once()

				_ = walkFn(fullPath, suite.mockDir, nil)
			}).Once()

			called := false
			err := suite.fs.Walk(suite.path, func(path string, info FileInfo) error {
				called = true
				suite.Equal(fullPath, path)
				suite.Equal(int64(100), info.Size)
				suite.Equal(tt.isDir, info.IsDir)
				return nil
			})
			suite.NoError(err)
			suite.True(called)
		})
	}
}

func (suite *FileSystemTestSuite) TestWalk_RelPathError() {
	root := "/absolute/root"
	incompatiblePath := "relative/path/file.jpg"
	simulatedErr := errors.New("rel: simulated error")

	suite.mockOS.On("Walk", root, mock.Anything).Return(simulatedErr).Run(func(args mock.Arguments) {
		walkFn := args.Get(1).(func(string, os.DirEntry, error) error)

		innerErr := walkFn(incompatiblePath, suite.mockDir, nil)
		suite.Error(innerErr)

		var errRel *ErrRelPath
		suite.True(errors.As(innerErr, &errRel), "internal error should be wrapped in ErrRelPath")
	}).Once()

	err := suite.fs.Walk(root, func(path string, info FileInfo) error {
		return nil
	})

	suite.Error(err)
	var errWalk *ErrWalk
	suite.True(errors.As(err, &errWalk), "final error should be wrapped in ErrWalk")
	suite.ErrorIs(err, simulatedErr)
	suite.EqualError(err, (&ErrWalk{Path: root, Err: simulatedErr}).Error())

	innerMockErr := &ErrRelPath{Root: root, Target: incompatiblePath, Err: simulatedErr}
	suite.Contains(innerMockErr.Error(), "failed to calculate relative path")
	suite.Equal(simulatedErr, innerMockErr.Unwrap())
}

func (suite *FileSystemTestSuite) TestWalk_InfoError() {
	expectedErr := errors.New("info error")

	suite.mockOS.On("Walk", suite.path, mock.Anything).Return(expectedErr).Run(func(args mock.Arguments) {
		walkFn := args.Get(1).(func(string, os.DirEntry, error) error)
		suite.mockDir.On("Info").Return(nil, expectedErr).Once()

		innerErr := walkFn(suite.path+"/image.jpg", suite.mockDir, nil)
		suite.Error(innerErr)

		var errInfo *ErrFileInfo
		suite.True(errors.As(innerErr, &errInfo), "internal error should be wrapped in ErrFileInfo")
	}).Once()

	err := suite.fs.Walk(suite.path, func(path string, info FileInfo) error {
		return nil
	})

	suite.Error(err)
	var errWalk *ErrWalk
	suite.True(errors.As(err, &errWalk), "final error should be wrapped in ErrWalk")
	suite.ErrorIs(err, expectedErr)
	suite.EqualError(err, (&ErrWalk{Path: suite.path, Err: expectedErr}).Error())

	innerMockErr := &ErrFileInfo{Path: suite.path + "/image.jpg", Err: expectedErr}
	suite.Contains(innerMockErr.Error(), "failed to get file info")
	suite.Equal(expectedErr, innerMockErr.Unwrap())
}

func (suite *FileSystemTestSuite) TestWalk_Success() {
	expectedPath := suite.path + "/trip/image.jpg"
	expectedRelPath := "trip/image.jpg"
	expectedSize := int64(1024)

	suite.mockOS.On("Walk", suite.path, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		walkFn := args.Get(1).(func(string, os.DirEntry, error) error)

		suite.mockDir.On("IsDir").Return(false).Once()
		suite.mockDir.On("Info").Return(FakeFileInfo{size: expectedSize}, nil).Once()

		_ = walkFn(expectedPath, suite.mockDir, nil)
	})

	var capturedInfo FileInfo
	err := suite.fs.Walk(suite.path, func(path string, info FileInfo) error {
		capturedInfo = info
		return nil
	})

	suite.NoError(err)
	suite.Equal(expectedPath, capturedInfo.Path)
	suite.Equal(expectedSize, capturedInfo.Size)
	suite.Equal(expectedRelPath, capturedInfo.RelPath)
}

func (suite *FileSystemTestSuite) TestErrors_Unwrap() {
	baseErr := errors.New("base error")

	tests := []struct {
		errWrapper interface{ Unwrap() error }
	}{
		{&ErrOpenDir{Err: baseErr}},
		{&ErrReadDir{Err: baseErr}},
		{&ErrCreateDir{Err: baseErr}},
		{&ErrCreateSiblingDir{Err: baseErr}},
		{&ErrReadFile{Err: baseErr}},
		{&ErrWriteFile{Err: baseErr}},
		{&ErrWalk{Err: baseErr}},
	}

	for _, tt := range tests {
		suite.Equal(baseErr, tt.errWrapper.Unwrap(), "Unwrap should return the underlying error")
	}
}

func TestFileSystemTestSuite(t *testing.T) {
	suite.Run(t, new(FileSystemTestSuite))
}
