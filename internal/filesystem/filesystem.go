package filesystem

import (
	"os"
	"path/filepath"

	"github.com/felipesimis/compactify-cli/internal/utils"
)

type FileSystem interface {
	ReadDir(path string) ([]string, error)
	CreateSiblingDir(path, suffix string) (string, error)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
}

type FileSystemWrapper struct {
	Mkdirer Mkdirer
}

type Dir interface {
	Readdir(count int) ([]os.FileInfo, error)
}

func NewFileSystem() FileSystem {
	return &FileSystemWrapper{Mkdirer: &OSWrapper{}}
}

func (fs *FileSystemWrapper) ReadDir(path string) ([]string, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, &ErrOpenDir{Err: err}
	}
	defer dir.Close()

	return fs.readDir(dir, path)
}

func (fs *FileSystemWrapper) readDir(dir Dir, path string) ([]string, error) {
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, &ErrReadDir{Path: path, Err: err}
	}

	var files []string
	for _, fi := range fileInfos {
		if !fi.IsDir() && utils.IsValidImage(fi.Name()) {
			files = append(files, filepath.Join(path, fi.Name()))
		}
	}
	return files, nil
}

type Mkdirer interface {
	Mkdir(name string, perm os.FileMode) error
}

type OSWrapper struct{}

func (o *OSWrapper) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (fs *FileSystemWrapper) CreateSiblingDir(path, suffix string) (string, error) {
	parentDir := filepath.Dir(path)
	newDir := filepath.Join(parentDir, filepath.Base(path)+suffix)
	err := fs.Mkdirer.Mkdir(newDir, os.ModePerm)
	if err != nil {
		return "", &ErrCreateSiblingDir{Err: err}
	}
	return newDir, nil
}

func (fs *FileSystemWrapper) ReadFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, &ErrReadFile{Path: path, Err: err}
	}
	return file, nil
}

func (fs *FileSystemWrapper) WriteFile(name string, data []byte) error {
	err := os.WriteFile(name, data, 0644)
	if err != nil {
		return &ErrWriteFile{Path: name, Err: err}
	}
	return nil
}
