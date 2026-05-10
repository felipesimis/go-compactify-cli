package processing

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/stretchr/testify/mock"
)

type mockFastFS struct {
	mock.Mock
	filesystem.FileSystem
	precomputedFiles []filesystem.FileInfo
}

func (m *mockFastFS) Walk(root string, walkFn func(path string, info filesystem.FileInfo) error) error {
	for _, info := range m.precomputedFiles {
		if err := walkFn(info.Path, info); err != nil {
			return err
		}
	}
	return nil
}

func (m *mockFastFS) CreateDir(name string) error {
	return nil
}

func BenchmarkDiscoverAndPrepare_Recursive_MassiveTree(b *testing.B) {
	mockFS := new(mockFastFS)

	for i := range 100 {
		dirRel := fmt.Sprintf("folder_%d", i)
		mockFS.precomputedFiles = append(mockFS.precomputedFiles, filesystem.FileInfo{
			Path:    filepath.Join("/fake/input", dirRel),
			RelPath: dirRel,
			IsDir:   true,
		})

		for j := range 100 {
			fileRel := filepath.Join(dirRel, fmt.Sprintf("img_%d.jpg", j))
			mockFS.precomputedFiles = append(mockFS.precomputedFiles, filesystem.FileInfo{
				Path:    filepath.Join("/fake/input", fileRel),
				RelPath: fileRel,
				IsDir:   false,
				Size:    2048 * 1024,
			})
		}
	}

	for b.Loop() {
		files, err := DiscoverAndPrepare(mockFS, "input", "output", true)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}

		if len(files) != 10000 {
			b.Fatalf("Expected 10,000 files, got %d", len(files))
		}
	}
}
