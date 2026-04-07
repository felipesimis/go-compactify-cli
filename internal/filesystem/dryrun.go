package filesystem

import (
	"io"
)

type DryRunFileSystem struct {
	original FileSystem
}

func NewDryRunFileSystem(original FileSystem) FileSystem {
	return &DryRunFileSystem{original: original}
}

func (d *DryRunFileSystem) CreateDir(name string) error {
	return nil
}

func (d *DryRunFileSystem) CreateSiblingDir(path, suffix string) (string, error) {
	return "", nil
}

func (d *DryRunFileSystem) ReadDir(path string) ([]FileInfo, error) {
	return d.original.ReadDir(path)
}

func (d *DryRunFileSystem) ReadFile(path string) ([]byte, error) {
	return d.original.ReadFile(path)
}

func (d *DryRunFileSystem) OpenFile(path string) (io.ReadCloser, error) {
	return d.original.OpenFile(path)
}

func (d *DryRunFileSystem) WriteFile(path string, data []byte) error {
	return nil
}
