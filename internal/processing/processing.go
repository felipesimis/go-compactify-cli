package processing

import (
	"log"
	"sync"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
)

type ProgressBarInterface interface {
	Increment()
	Finish()
}

type FileProcessingParams struct {
	File        filesystem.FileInfo
	FS          filesystem.FileSystem
	OutputDir   string
	ProgressBar ProgressBarInterface
	ExtraParams interface{}
}

type fileProcessorFunc func(FileProcessingParams) error

type ProcessFilesParams struct {
	Files         []filesystem.FileInfo
	FS            filesystem.FileSystem
	OutputDir     string
	ProgressBar   ProgressBarInterface
	ExtraParams   interface{}
	ProcessorFunc fileProcessorFunc
	Concurrency   int
}

func ProcessFiles(params ProcessFilesParams) {
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
				errChan <- err
			}
			fpParams.ProgressBar.Increment()
		}(file)
	}

	wg.Wait()
	close(sem)
	close(errChan)

	for err := range errChan {
		if err != nil {
			log.Println(err)
		}
	}
}
