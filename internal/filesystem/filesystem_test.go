package filesystem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockDir struct {
	mock.Mock
}

func (m *MockDir) Readdir(count int) ([]os.FileInfo, error) {
	args := m.Called(count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

type MockMkdirer struct {
	Err error
}

func (m *MockMkdirer) Mkdir(name string, perm os.FileMode) error {
	return m.Err
}

type FileSystemTestSuite struct {
	suite.Suite
	fs     FileSystem
	tmpDir string
}

func (suite *FileSystemTestSuite) SetupTest() {
	suite.fs = NewFileSystem()
	var err error
	suite.tmpDir, err = os.MkdirTemp("", "test")
	assert.NoError(suite.T(), err)
}

func (suite *FileSystemTestSuite) TearDownTest() {
	defer os.RemoveAll(suite.tmpDir)
}

func (suite *FileSystemTestSuite) TestReadDir() {
	imageFiles := []string{"image1.jpg", "image2.jpeg", "image3.png", "image4.webp"}
	files := append(imageFiles, "file1.txt", "file2.pdf", "file3.doc")
	for _, file := range files {
		tmpFile, err := os.Create(filepath.Join(suite.tmpDir, file))
		assert.NoError(suite.T(), err)
		tmpFile.Close()
	}

	result, err := suite.fs.ReadDir(suite.tmpDir)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 4)

	for _, file := range result {
		assert.Contains(suite.T(), files, filepath.Base(file.Path))
	}
}

func (suite *FileSystemTestSuite) TestOpenError() {
	nonExistentPath := filepath.Join(suite.tmpDir, "nonexistent")
	files, err := suite.fs.ReadDir(nonExistentPath)
	expectedErr := &ErrOpenDir{Err: fmt.Errorf("open %s: no such file or directory", nonExistentPath)}
	assert.Nil(suite.T(), files)
	assert.EqualError(suite.T(), err, expectedErr.Error())
}

func (suite *FileSystemTestSuite) TestReaddirError() {
	mockDir := new(MockDir)
	mockDir.On("Readdir", -1).Return(nil, errors.New("simulated readdir error"))

	files, err := suite.fs.(*FileSystemWrapper).readDir(mockDir, "/mock/path")
	expectedErr := &ErrReadDir{Path: "/mock/path", Err: errors.New("simulated readdir error")}
	assert.Nil(suite.T(), files)
	assert.EqualError(suite.T(), err, expectedErr.Error())

	mockDir.AssertExpectations(suite.T())
}

func TestFileSystemWrapper_CreateSiblingDirError(t *testing.T) {
	mockMkdirer := &MockMkdirer{Err: errors.New("mock error")}
	fs := &FileSystemWrapper{Mkdirer: mockMkdirer}

	newDir, err := fs.CreateSiblingDir("/some/path", "_suffix")
	exportedErr := &ErrCreateSiblingDir{Err: mockMkdirer.Err}
	assert.Empty(t, newDir)
	assert.EqualError(t, err, exportedErr.Error())
}

func (suite *FileSystemTestSuite) TestCreateSiblingDir() {
	newDir, err := suite.fs.CreateSiblingDir(suite.tmpDir, "_suffix")
	assert.NoError(suite.T(), err)
	assert.DirExists(suite.T(), newDir)
	assert.Contains(suite.T(), newDir, "_suffix")
}

func (suite *FileSystemTestSuite) TestReadFileError() {
	nonExistentPath := filepath.Join(suite.tmpDir, "nonexistent")
	data, err := suite.fs.ReadFile(nonExistentPath)
	expectedErr := &ErrReadFile{Path: nonExistentPath, Err: fmt.Errorf("open %s: no such file or directory", nonExistentPath)}
	assert.Nil(suite.T(), data)
	assert.EqualError(suite.T(), err, expectedErr.Error())
}

func (suite *FileSystemTestSuite) TestReadFile() {
	tmpFile, err := os.CreateTemp(suite.tmpDir, "file*.txt")
	assert.NoError(suite.T(), err)
	defer os.Remove(tmpFile.Name())

	data, err := suite.fs.ReadFile(tmpFile.Name())
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), data)
}

func (suite *FileSystemTestSuite) TestWriteFileError() {
	invalidPath := "/invalid/path/to/image.png"
	err := suite.fs.WriteFile(invalidPath, []byte("data"))
	expectedErr := &ErrWriteFile{Path: invalidPath, Err: fmt.Errorf("open %s: no such file or directory", invalidPath)}
	assert.EqualError(suite.T(), err, expectedErr.Error())
}

func (suite *FileSystemTestSuite) TestWriteFile() {
	tmpFile := filepath.Join(suite.tmpDir, "image.png")
	err := suite.fs.WriteFile(tmpFile, []byte("data"))
	assert.NoError(suite.T(), err)
	assert.FileExists(suite.T(), tmpFile)
}

func TestFileSystemTestSuite(t *testing.T) {
	suite.Run(t, new(FileSystemTestSuite))
}
