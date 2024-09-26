package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	expectedErr := &ErrOpenDir{Path: nonExistentPath, Err: fmt.Errorf("open %s: no such file or directory", nonExistentPath)}
	assert.Nil(t, files)
	assert.EqualError(t, err, expectedErr.Error())
}
