package filesystem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestFileSystemWrapper_ReadDir(t *testing.T) {
	fs := NewFileSystem()

	tmpDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	imageFiles := []string{"image1.jpg", "image2.jpeg", "image3.png", "image4.webp"}
	files := append(imageFiles, "file1.txt", "file2.pdf", "file3.doc")
	for _, file := range files {
		tmpFile, err := os.Create(filepath.Join(tmpDir, file))
		assert.NoError(t, err)
		tmpFile.Close()
	}

	result, err := fs.ReadDir(tmpDir)
	assert.NoError(t, err)
	assert.Len(t, result, 4)

	for _, file := range result {
		assert.Contains(t, files, filepath.Base(file))
	}
}

func TestFileSystemWrapper_OpenError(t *testing.T) {
	fs := NewFileSystem()

	tmpDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)

	nonExistentPath := filepath.Join(tmpDir, "nonexistent")
	files, err := fs.ReadDir(nonExistentPath)
	expectedErr := &ErrOpenDir{Err: fmt.Errorf("open %s: no such file or directory", nonExistentPath)}
	assert.Nil(t, files)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestFileSystemWrapper_ReaddirError(t *testing.T) {
	fs := NewFileSystem()

	mockDir := new(MockDir)
	mockDir.On("Readdir", -1).Return(nil, errors.New("simulated readdir error"))

	files, err := fs.(*FileSystemWrapper).readDir(mockDir, "/mock/path")
	expectedErr := &ErrReadDir{Path: "/mock/path", Err: errors.New("simulated readdir error")}
	assert.Nil(t, files)
	assert.EqualError(t, err, expectedErr.Error())

	mockDir.AssertExpectations(t)
}

func TestFileSystemWrapper_CreateSiblingDirError(t *testing.T) {
	mockMkdirer := &MockMkdirer{Err: errors.New("mock error")}
	fs := &FileSystemWrapper{Mkdirer: mockMkdirer}

	newDir, err := fs.CreateSiblingDir("/some/path", "_suffix")
	exportedErr := &ErrCreateSiblingDir{Err: mockMkdirer.Err}
	assert.Empty(t, newDir)
	assert.EqualError(t, err, exportedErr.Error())
}

func TestFileSystemWrapper_CreateSiblingDir(t *testing.T) {
	fs := NewFileSystem()

	tmpDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	newDir, err := fs.CreateSiblingDir(tmpDir, "_suffix")
	assert.NoError(t, err)
	assert.DirExists(t, newDir)
	assert.Contains(t, newDir, "_suffix")
}

func TestFileSystemWrapper_ReadFileError(t *testing.T) {
	fs := NewFileSystem()
	tmpDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)

	nonExistentPath := filepath.Join(tmpDir, "nonexistent")
	data, err := fs.ReadFile(nonExistentPath)
	expectedErr := &ErrReadFile{Path: nonExistentPath, Err: fmt.Errorf("open %s: no such file or directory", nonExistentPath)}
	assert.Nil(t, data)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestFileSystemWrapper_ReadFile(t *testing.T) {
	fs := NewFileSystem()
	tmpDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	tmpFile, err := os.CreateTemp(tmpDir, "file*.txt")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	data, err := fs.ReadFile(tmpFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, data)
}

func TestFileSystemWrapper_WriteFileError(t *testing.T) {
	fs := NewFileSystem()

	invalidPath := "/invalid/path/to/image.png"
	err := fs.WriteFile(invalidPath, []byte("data"))
	expectedErr := &ErrWriteFile{Path: invalidPath, Err: fmt.Errorf("open %s: no such file or directory", invalidPath)}
	assert.EqualError(t, err, expectedErr.Error())
}
