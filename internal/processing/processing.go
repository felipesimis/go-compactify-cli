package processing

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
)

type ProgressBarInterface interface {
	Increment()
	Finish()
}

type ProcessorFS interface {
	filesystem.FileReader
	filesystem.FileWriter
	CreateDir(path string) error
}

type FileTask struct {
	File        filesystem.FileInfo
	FS          ProcessorFS
	InputDir    string
	OutputDir   string
	ProgressBar ProgressBarInterface
	ExtraParams interface{}
}

type fileTaskHandler func(task FileTask) error

type FileBatchConfig struct {
	Files       []filesystem.FileInfo
	FS          ProcessorFS
	InputDir    string
	OutputDir   string
	ProgressBar ProgressBarInterface
	ExtraParams interface{}
	Handler     fileTaskHandler
	Concurrency int
}

func RunFileBatch(config FileBatchConfig) []error {
	concurrency := config.Concurrency
	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
	}

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	errChan := make(chan error, concurrency)

	var processErrors []error
	errDone := make(chan struct{})
	go func() {
		for err := range errChan {
			processErrors = append(processErrors, err)
		}
		close(errDone)
	}()

	for _, file := range config.Files {
		wg.Add(1)
		sem <- struct{}{}
		go func(file filesystem.FileInfo) {
			defer func() {
				<-sem
				wg.Done()
			}()
			err := config.Handler(FileTask{
				File:        file,
				FS:          config.FS,
				InputDir:    config.InputDir,
				OutputDir:   config.OutputDir,
				ProgressBar: config.ProgressBar,
				ExtraParams: config.ExtraParams,
			})
			config.ProgressBar.Increment()

			if err != nil {
				errChan <- fmt.Errorf("error processing file '%s': %w", file.Path, err)
			}
		}(file)
	}

	wg.Wait()
	close(sem)
	close(errChan)
	<-errDone

	return processErrors
}
