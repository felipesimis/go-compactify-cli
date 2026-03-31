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
	sem := make(chan struct{}, params.Concurrency)
	var wg sync.WaitGroup
	var errChan = make(chan error, len(params.Files))

	for _, file := range params.Files {
		wg.Add(1)
		sem <- struct{}{}
		go func(file filesystem.FileInfo) {
			defer func() {
				<-sem
				wg.Done()
			}()
			fpParams := FileProcessingParams{
				File:        file,
				FS:          params.FS,
				OutputDir:   params.OutputDir,
				ProgressBar: params.ProgressBar,
				ExtraParams: params.ExtraParams,
			}
			err := params.ProcessorFunc(fpParams)
			if err != nil {
				errChan <- fmt.Errorf("error processing file '%s': %w", file.Path, err)
			}
			fpParams.ProgressBar.Increment()
		}(file)
	}

	wg.Wait()
	close(sem)
	close(errChan)

	var processErrors []error

	for err := range errChan {
		if err != nil {
			processErrors = append(processErrors, err)
		}
	}
	return processErrors
}
