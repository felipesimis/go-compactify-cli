package processing

import (
	"path/filepath"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
)

func DiscoverAndPrepare(fs filesystem.FileSystem, inputDir, outputDir string, recursive bool) ([]filesystem.FileInfo, error) {
	filesToProcess := make([]filesystem.FileInfo, 0, 1024)

	if recursive {
		err := fs.Walk(inputDir, func(path string, info filesystem.FileInfo) error {
			if info.IsDir {
				if info.RelPath == "." {
					return nil
				}
				destDir := filepath.Join(outputDir, info.RelPath)
				if err := fs.CreateDir(destDir); err != nil {
					return err
				}
				return nil
			}

			if utils.IsValidImage(info.RelPath) {
				filesToProcess = append(filesToProcess, info)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		entries, err := fs.ReadDir(inputDir)
		if err != nil {
			return nil, err
		}

		for _, info := range entries {
			if !info.IsDir && utils.IsValidImage(info.RelPath) {
				filesToProcess = append(filesToProcess, info)
			}
		}
	}

	return filesToProcess, nil
}
