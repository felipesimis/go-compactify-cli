package filesystem

import (
	"io"
	"os"
	"path/filepath"

	"github.com/felipesimis/compactify-cli/internal/utils"
)

type FileReader interface {
	ReadFile(path string) ([]byte, error)
	OpenFile(path string) (io.ReadCloser, error)
}

type FileWriter interface {
	WriteFile(path string, data []byte) error
}

type FileSystem interface {
	ReadDir(path string) ([]FileInfo, error)
	CreateDir(name string) error
	CreateSiblingDir(path, suffix string) (string, error)
	FileReader
	FileWriter
}

type FileInfo struct {
	Path string
	Size int64
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

func (fs *FileSystemWrapper) ReadDir(path string) ([]FileInfo, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, &ErrOpenDir{Err: err}
	}
	defer dir.Close()

	return fs.readDir(dir, path)
}

func (fs *FileSystemWrapper) readDir(dir Dir, path string) ([]FileInfo, error) {
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, &ErrReadDir{Path: path, Err: err}
	}

	var files []FileInfo
	for _, fi := range fileInfos {
		if !fi.IsDir() && utils.IsValidImage(fi.Name()) {
			files = append(files, FileInfo{
				Path: filepath.Join(path, fi.Name()),
				Size: fi.Size(),
			})
		}
	}
	return files, nil
}

type Mkdirer interface {
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
}

type OSWrapper struct{}

func (o *OSWrapper) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (o *OSWrapper) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs *FileSystemWrapper) CreateDir(name string) error {
	err := fs.Mkdirer.MkdirAll(name, os.ModePerm)
	if err != nil {
		return &ErrCreateDir{Path: name, Err: err}
	}
	return nil
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

func (fs *FileSystemWrapper) OpenFile(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
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
