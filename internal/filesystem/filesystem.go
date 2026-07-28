package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"strings"
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
	Walk(root string, walkFn func(path string, info FileInfo) error) error
	CreateDir(name string) error
	CreateSiblingDir(path, suffix string) (string, error)
	Stat(path string) (FileInfo, error)
	FileReader
	FileWriter
}

type File interface {
	io.ReadCloser
	Readdir(count int) ([]os.FileInfo, error)
}

type FileInfo struct {
	Path    string
	RelPath string
	Size    int64
	IsDir   bool
}

type OSOperations interface {
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Open(name string) (File, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	Walk(root string, walkFn func(path string, d os.DirEntry, err error) error) error
	Stat(name string) (os.FileInfo, error)
	Rename(oldPath, newPath string) error
	Remove(name string) error
}

type FileSystemWrapper struct {
	os OSOperations
}

type Dir interface {
	Readdir(count int) ([]os.FileInfo, error)
}

func NewFileSystem() FileSystem {
	return &FileSystemWrapper{os: &OSWrapper{}}
}

func (fs *FileSystemWrapper) ReadDir(path string) ([]FileInfo, error) {
	dir, err := fs.os.Open(path)
	if err != nil {
		return nil, &ErrOpenDir{Path: path, Err: err}
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
		files = append(files, FileInfo{
			Path:    filepath.Join(path, fi.Name()),
			RelPath: fi.Name(),
			Size:    fi.Size(),
			IsDir:   fi.IsDir(),
		})
	}
	return files, nil
}

type OSWrapper struct{}

func (o *OSWrapper) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (o *OSWrapper) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (o *OSWrapper) Open(name string) (File, error) {
	return os.Open(name)
}

func (o *OSWrapper) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (o *OSWrapper) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (o *OSWrapper) Walk(root string, walkFn func(path string, d os.DirEntry, err error) error) error {
	return filepath.WalkDir(root, walkFn)
}

func (o *OSWrapper) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (o *OSWrapper) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func (o *OSWrapper) Remove(name string) error {
	return os.Remove(name)
}

func (fs *FileSystemWrapper) CreateDir(name string) error {
	if err := fs.os.MkdirAll(name, os.ModePerm); err != nil {
		return &ErrCreateDir{Path: name, Err: err}
	}
	return nil
}

func (fs *FileSystemWrapper) CreateSiblingDir(path, suffix string) (string, error) {
	parentDir := filepath.Dir(path)
	newDir := filepath.Join(parentDir, filepath.Base(path)+suffix)
	if err := fs.os.MkdirAll(newDir, os.ModePerm); err != nil {
		return "", &ErrCreateSiblingDir{Path: path, Err: err}
	}
	return newDir, nil
}

func (fs *FileSystemWrapper) ReadFile(path string) ([]byte, error) {
	file, err := fs.os.ReadFile(path)
	if err != nil {
		return nil, &ErrReadFile{Path: path, Err: err}
	}
	return file, nil
}

func (fs *FileSystemWrapper) OpenFile(path string) (io.ReadCloser, error) {
	file, err := fs.os.Open(path)
	if err != nil {
		return nil, &ErrReadFile{Path: path, Err: err}
	}
	return file, nil
}

func (fs *FileSystemWrapper) WriteFile(name string, data []byte) error {
	tmpName := name + ".tmp"

	if err := fs.os.WriteFile(tmpName, data, 0644); err != nil {
		return &ErrWriteFile{Path: tmpName, Err: err}
	}

	if err := fs.os.Rename(tmpName, name); err != nil {
		_ = fs.os.Remove(tmpName)
		return &ErrWriteFile{Path: name, Err: err}
	}
	return nil
}

func (fs *FileSystemWrapper) Stat(path string) (FileInfo, error) {
	file, err := fs.os.Stat(path)
	if err != nil {
		return FileInfo{}, &ErrFileInfo{Path: path, Err: err}
	}
	return FileInfo{
		Path:    path,
		Size:    file.Size(),
		RelPath: filepath.Base(path),
		IsDir:   file.IsDir(),
	}, nil
}

func (fs *FileSystemWrapper) Walk(root string, walkFn func(path string, info FileInfo) error) error {
	cleanRoot := filepath.Clean(root)

	err := fs.os.Walk(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath := path
		if cleanRoot == "." {
			relPath = path
		} else if strings.HasPrefix(path, cleanRoot) {
			relPath = path[len(cleanRoot):]
			if len(relPath) > 0 && os.IsPathSeparator(relPath[0]) {
				relPath = relPath[1:]
			}
		}
		if relPath == "" {
			relPath = "."
		}

		info, err := d.Info()
		if err != nil {
			return &ErrFileInfo{Path: path, Err: err}
		}

		return walkFn(path, FileInfo{
			Path:    path,
			RelPath: relPath,
			Size:    info.Size(),
			IsDir:   d.IsDir(),
		})
	})

	if err != nil {
		return &ErrWalk{Path: root, Err: err}
	}
	return nil
}
