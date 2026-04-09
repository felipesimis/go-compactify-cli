package filesystem

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *FileSystemTestSuite) SetupTest() {
	suite.mockOS = new(MockOSOperations)
	suite.fs = &FileSystemWrapper{os: suite.mockOS}
	suite.mockFile = new(MockFile)
	suite.path = "/mock/dir"
}

func (suite *FileSystemTestSuite) TestNewFileSystem() {
	fs := NewFileSystem()
	suite.NotNil(fs)
	_, ok := fs.(*FileSystemWrapper)
	suite.True(ok)
}

func (suite *FileSystemTestSuite) TestReadDir() {
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
	suite.Len(result, 4)
	for _, file := range result {
		suite.Contains([]string{"image1.jpg", "image2.jpeg", "image3.png", "image4.webp"}, filepath.Base(file.Path))
	}
	suite.mockOS.AssertExpectations(suite.T())
	suite.mockFile.AssertExpectations(suite.T())
}

func (suite *FileSystemTestSuite) TestReaddir_OpenError() {
	suite.mockOS.On("Open", suite.path).Return(nil, errors.New("simulated open error"))

	result, err := suite.fs.ReadDir(suite.path)
	expectedErr := &ErrOpenDir{Err: errors.New("simulated open error")}
	suite.Nil(result)
	suite.EqualError(err, expectedErr.Error())
	suite.mockOS.AssertExpectations(suite.T())
}

func (suite *FileSystemTestSuite) TestReaddir_ReadDirError() {
	suite.mockOS.On("Open", suite.path).Return(suite.mockFile, nil)
	suite.mockFile.On("Readdir", -1).Return(nil, errors.New("simulated readdir error"))
	suite.mockFile.On("Close").Return(nil)

	result, err := suite.fs.ReadDir(suite.path)
	expectedErr := &ErrReadDir{Path: suite.path, Err: errors.New("simulated readdir error")}
	suite.Nil(result)
	suite.EqualError(err, expectedErr.Error())
	suite.mockOS.AssertExpectations(suite.T())
	suite.mockFile.AssertExpectations(suite.T())
}

func (suite *FileSystemTestSuite) TestFileSystemWrapper_CreateDirError() {
	suite.mockOS.On("MkdirAll", suite.path, os.ModePerm).Return(errors.New("mock error"))

	err := suite.fs.CreateDir(suite.path)
	expectedErr := &ErrCreateDir{Path: suite.path, Err: errors.New("mock error")}
	suite.EqualError(err, expectedErr.Error())
}

func (suite *FileSystemTestSuite) TestFileSystemWrapper_CreateDir() {
	suite.mockOS.On("MkdirAll", suite.path, os.ModePerm).Return(nil)

	err := suite.fs.CreateDir(suite.path)
	suite.NoError(err)
	suite.mockOS.AssertExpectations(suite.T())
}

func (suite *FileSystemTestSuite) TestFileSystemWrapper_CreateSiblingDirError() {
	expectedPath := suite.path + "-suffix"
	expectedErr := &ErrCreateSiblingDir{Err: errors.New("mock error")}
	suite.mockOS.On("Mkdir", expectedPath, os.ModePerm).Return(errors.New("mock error"))

	newDir, err := suite.fs.CreateSiblingDir(suite.path, "-suffix")
	suite.Empty(newDir)
	suite.EqualError(err, expectedErr.Error())
	suite.mockOS.AssertExpectations(suite.T())
}

func (suite *FileSystemTestSuite) TestFileSystemWrapper_CreateSiblingDir() {
	expectedPath := suite.path + "-suffix"
	suite.mockOS.On("Mkdir", expectedPath, os.ModePerm).Return(nil)

	newDir, err := suite.fs.CreateSiblingDir(suite.path, "-suffix")
	suite.NoError(err)
	suite.Equal(expectedPath, newDir)
	suite.mockOS.AssertExpectations(suite.T())
}

func (suite *FileSystemTestSuite) TestReadFileError() {
	expectedErr := &ErrReadFile{Path: suite.path, Err: errors.New("mock error")}
	suite.mockOS.On("ReadFile", suite.path).Return(nil, expectedErr.Err)

	data, err := suite.fs.ReadFile(suite.path)
	suite.Nil(data)
	suite.EqualError(err, expectedErr.Error())
	suite.mockOS.AssertExpectations(suite.T())
}

func (suite *FileSystemTestSuite) TestReadFile() {
	expectedData := []byte("file content")
	suite.mockOS.On("ReadFile", suite.path).Return(expectedData, nil)

	data, err := suite.fs.ReadFile(suite.path)
	suite.NoError(err)
	suite.Equal(expectedData, data)
	suite.mockOS.AssertExpectations(suite.T())
}

func (suite *FileSystemTestSuite) TestOpenFileError() {
	expectedErr := &ErrReadFile{Path: suite.path, Err: errors.New("mock error")}
	suite.mockOS.On("Open", suite.path).Return(nil, expectedErr.Err)

	file, err := suite.fs.OpenFile(suite.path)
	suite.Nil(file)
	suite.EqualError(err, expectedErr.Error())
	suite.mockOS.AssertExpectations(suite.T())
}

func (suite *FileSystemTestSuite) TestOpenFile() {
	suite.mockOS.On("Open", suite.path).Return(suite.mockFile, nil)

	file, err := suite.fs.OpenFile(suite.path)
	suite.NoError(err)
	suite.Equal(suite.mockFile, file)
	suite.mockOS.AssertExpectations(suite.T())
}

func (suite *FileSystemTestSuite) TestWriteFileError() {
	expectedErr := &ErrWriteFile{Path: suite.path, Err: errors.New("mock error")}
	suite.mockOS.On("WriteFile", suite.path, []byte("data"), os.FileMode(0644)).Return(expectedErr.Err)

	err := suite.fs.WriteFile(suite.path, []byte("data"))
	suite.EqualError(err, expectedErr.Error())
	suite.mockOS.AssertExpectations(suite.T())
}

func (suite *FileSystemTestSuite) TestWriteFile() {
	data := []byte("data")
	suite.mockOS.On("WriteFile", suite.path, data, os.FileMode(0644)).Return(nil)

	err := suite.fs.WriteFile(suite.path, data)
	suite.NoError(err)
	suite.mockOS.AssertExpectations(suite.T())
}

func TestFileSystemTestSuite(t *testing.T) {
	suite.Run(t, new(FileSystemTestSuite))
}
