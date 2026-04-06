package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/processing"
	"github.com/felipesimis/compactify-cli/internal/utils"
	"github.com/felipesimis/compactify-cli/pkg/progress"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type OperationConfig struct {
	Ctx                context.Context
	FileSystem         filesystem.FileSystem
	InputDir           string
	OutputSuffix       string
	ProgressBarMessage string
	ExtraParams        interface{}
	ProcessorFunc      func(ctx context.Context, p processing.FileProcessingParams, stats *utils.ImageProcessingStats) error
	ResultVerb         string
}

func RunOperation(config OperationConfig) {
	files, err := config.FileSystem.ReadDir(config.InputDir)
	if err != nil {
		log.Fatal(err)
	}

	outputDir, err := config.FileSystem.CreateSiblingDir(config.InputDir, config.OutputSuffix)
	if err != nil {
		log.Fatal(err)
	}

	stats := &utils.ImageProcessingStats{}
	resultBuilder := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(os.Stdout, len(files), concurrency, config.ProgressBarMessage)

	wrappedProcessor := func(p processing.FileProcessingParams) error {
		return config.ProcessorFunc(config.Ctx, p, stats)
	}

	params := processing.ProcessFilesParams{
		Files:         files,
		FS:            config.FileSystem,
		OutputDir:     outputDir,
		ProgressBar:   progressBar,
		ExtraParams:   config.ExtraParams,
		ProcessorFunc: wrappedProcessor,
		Concurrency:   concurrency,
	}
	processErrors := processing.ProcessFiles(params)

	progressBar.Finish()

	totalImages := uint32(len(files))
	resultBuilder.SetTotalImages(totalImages).
		SetSkippedImages(stats.SkippedImages.Load()).
		SetProcessedImages(stats.ProcessedImages.Load()).
		SetOutputDirectory(outputDir).
		SetInitialSize(float64(stats.InitialSize.Load())).
		SetFinalSize(float64(stats.FinalSize.Load())).
		SetErrors(processErrors)
	result := resultBuilder.Build()
	fmt.Println(result.PrintResults(config.ResultVerb))
}

func HandleImageProcessing(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats, processFunc func([]byte) ([]byte, error)) error {
	select {
	case <-ctx.Done():
		stats.SkippedImages.Add(1)
		return ctx.Err()
	default:
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.Grow(int(params.File.Size))
	defer bufferPool.Put(buf)

	file, err := params.FS.OpenFile(params.File.Path)
	if err != nil {
		stats.SkippedImages.Add(1)
		return err
	}
	defer file.Close()

	_, err = io.Copy(buf, file)
	if err != nil {
		stats.SkippedImages.Add(1)
		return err
	}

	imgBytes := buf.Bytes()
	stats.InitialSize.Add(uint64(len(imgBytes)))

	newImg, err := processFunc(imgBytes)
	if err != nil {
		stats.SkippedImages.Add(1)
		return err
	}

	outputPath := determineOutputPath(params)
	err = params.FS.WriteFile(outputPath, newImg)
	if err != nil {
		stats.SkippedImages.Add(1)
		return err
	}

	stats.FinalSize.Add(uint64(len(newImg)))
	stats.ProcessedImages.Add(1)
	return nil
}

func determineOutputPath(params processing.FileProcessingParams) string {
	if convertParams, ok := params.ExtraParams.(ConvertParams); ok && convertParams.Format != "" {
		originalFileName := filepath.Base(params.File.Path)
		fileExt := filepath.Ext(originalFileName)
		fileNameWithoutExt := strings.TrimSuffix(originalFileName, fileExt)
		newFilename := fmt.Sprintf("%s.%s", fileNameWithoutExt, convertParams.Format)

		return filepath.Join(params.OutputDir, newFilename)
	}

	return utils.BuildOutputPath(params.OutputDir, params.File.Path)
}
