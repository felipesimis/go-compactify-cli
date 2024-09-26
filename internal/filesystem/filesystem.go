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

func NewFileSystem() FileSystem {
	return &FileSystemWrapper{}
}

func (fs *FileSystemWrapper) ReadDir(path string) ([]string, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, &ErrOpenDir{Path: path, Err: err}
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, fi := range fileInfos {
		if !fi.IsDir() && utils.IsValidImage(fi.Name()) {
			files = append(files, filepath.Join(path, fi.Name()))
		}
	}
	return files, nil
}
