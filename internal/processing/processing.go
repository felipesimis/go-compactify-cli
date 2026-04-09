package processing

import (
	"fmt"
	"sync"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
)

type ProgressBarInterface interface {
	Increment()
	Finish()
}

type FileProcessingParams struct {
	File        filesystem.FileInfo
	FS          FileReaderWriter
	OutputDir   string
	ProgressBar ProgressBarInterface
	ExtraParams interface{}
}

type fileProcessorFunc func(FileProcessingParams) error

type FileReaderWriter interface {
	filesystem.FileReader
	filesystem.FileWriter
}

type ProcessFilesParams struct {
	Files         []filesystem.FileInfo
	FS            FileReaderWriter
	OutputDir     string
	ProgressBar   ProgressBarInterface
	ExtraParams   interface{}
	ProcessorFunc fileProcessorFunc
	Concurrency   int
}

func ProcessFiles(params ProcessFilesParams) []error {
	concurrency := params.Concurrency
	if concurrency <= 0 {
		concurrency = 1
	}

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	errChan := make(chan error, len(params.Files))

	for _, file := range params.Files {
		wg.Add(1)
		sem <- struct{}{}
		go func(file filesystem.FileInfo) {
			defer func() {
				<-sem
				wg.Done()
			}()
			err := params.ProcessorFunc(FileProcessingParams{
				File:        file,
				FS:          params.FS,
				OutputDir:   params.OutputDir,
				ProgressBar: params.ProgressBar,
				ExtraParams: params.ExtraParams,
			})
			params.ProgressBar.Increment()

			if err != nil {
				errChan <- fmt.Errorf("error processing file '%s': %w", file.Path, err)
			}
		}(file)
	}

	wg.Wait()
	close(sem)
	close(errChan)

	var processErrors []error
	for err := range errChan {
		processErrors = append(processErrors, err)
	}
	return processErrors
}
