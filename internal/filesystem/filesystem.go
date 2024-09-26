package filesystem

import (
	"os"
	"path/filepath"

	"github.com/felipesimis/compactify-cli/internal/utils"
)

type FileSystem interface {
	ReadDir(path string) ([]string, error)
}

type FileSystemWrapper struct{}

type Dir interface {
	Readdir(count int) ([]os.FileInfo, error)
}

func NewFileSystem() FileSystem {
	return &FileSystemWrapper{}
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
		return nil, &ReadDirError{Path: path, Err: err}
	}

	var files []string
	for _, fi := range fileInfos {
		if !fi.IsDir() && utils.IsValidImage(fi.Name()) {
			files = append(files, filepath.Join(path, fi.Name()))
		}
	}
	return files, nil
}
