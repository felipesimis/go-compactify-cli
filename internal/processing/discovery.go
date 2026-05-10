package processing

import (
	"path/filepath"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
)

var filepathAbs = filepath.Abs

func DiscoverAndPrepare(fs filesystem.FileSystem, inputDir, outputDir string, recursive bool) ([]filesystem.FileInfo, error) {
	filesToProcess := make([]filesystem.FileInfo, 0, 1024)

	if recursive {
		absOutput, err := filepathAbs(outputDir)
		if err != nil {
			return nil, err
		}

		err = fs.Walk(inputDir, func(path string, info filesystem.FileInfo) error {
			if info.IsDir {
				absPath, err := filepathAbs(path)
				if err != nil {
					return err
				}
				if absPath == absOutput {
					return filepath.SkipDir
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
