package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"
)

type fakeDirEntry struct {
	info os.FileInfo
}

func (d fakeDirEntry) Name() string               { return d.info.Name() }
func (d fakeDirEntry) IsDir() bool                { return d.info.IsDir() }
func (d fakeDirEntry) Type() os.FileMode          { return d.info.Mode().Type() }
func (d fakeDirEntry) Info() (os.FileInfo, error) { return d.info, nil }

func BenchmarkFileSystem_PathResolution_Battle(b *testing.B) {
	root := "/fake/input"
	var paths []string

	for i := range 10000 {
		paths = append(paths, filepath.Join(root, fmt.Sprintf("folder_%d", i%100), fmt.Sprintf("img_%d.jpg", i)))
	}

	mockOS := new(MockOSOperations)

	mockOS.On("Walk", root, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		walkFn := args.Get(1).(func(string, os.DirEntry, error) error)
		dummyInfo := &FakeFileInfo{name: "img.jpg", size: 1024, isDir: false}
		dummyEntry := fakeDirEntry{info: dummyInfo}

		for _, p := range paths {
			_ = walkFn(p, dummyEntry, nil)
		}
	})

	fsWrapper := &FileSystemWrapper{os: mockOS}

	for b.Loop() {
		err := fsWrapper.Walk(root, func(path string, info FileInfo) error {
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
