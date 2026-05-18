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
}

type fileTaskHandler func(task FileTask) error

type FileBatchConfig struct {
	Files       []filesystem.FileInfo
	FS          ProcessorFS
	InputDir    string
	OutputDir   string
	ProgressBar ProgressBarInterface
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
	go func(ch <-chan error) {
		for err := range ch {
			processErrors = append(processErrors, err)
		}
		close(errDone)
	}(errChan)

	for _, file := range config.Files {
		wg.Add(1)
		sem <- struct{}{}
		go func(file filesystem.FileInfo, sendErr chan<- error) {
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
			})
			config.ProgressBar.Increment()

			if err != nil {
				sendErr <- fmt.Errorf("error processing file '%s': %w", file.Path, err)
			}
		}(file, errChan)
	}

	wg.Wait()
	close(sem)
	close(errChan)
	<-errDone

	return processErrors
}
