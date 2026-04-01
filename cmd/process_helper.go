package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/processing"
	"github.com/felipesimis/compactify-cli/internal/utils"
	"github.com/felipesimis/compactify-cli/pkg/progress"
)

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

	var initialSize, finalSize uint64
	var skippedImages, processedImages uint32

	resultBuilder := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(os.Stdout, len(files), concurrency, config.ProgressBarMessage)

	wrappedProcessor := func(p processing.FileProcessingParams) error {
		stats := &utils.ImageProcessingStats{
			InitialSize:     &initialSize,
			FinalSize:       &finalSize,
			SkippedImages:   &skippedImages,
			ProcessedImages: &processedImages,
		}
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
		SetSkippedImages(skippedImages).
		SetProcessedImages(processedImages).
		SetOutputDirectory(outputDir).
		SetInitialSize(float64(initialSize)).
		SetFinalSize(float64(finalSize)).
		SetErrors(processErrors)
	result := resultBuilder.Build()
	fmt.Println(result.PrintResults(config.ResultVerb))
}

func HandleImageProcessing(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats, processFunc func([]byte) ([]byte, error)) error {
	select {
	case <-ctx.Done():
		atomic.AddUint32(stats.SkippedImages, 1)
		return ctx.Err()
	default:
	}

	img, err := params.FS.ReadFile(params.File.Path)
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.InitialSize, uint64(params.File.Size))
	newImg, err := processFunc(img)
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	var outputPath string
	if convertParams, ok := params.ExtraParams.(ConvertParams); ok && convertParams.Format != "" {
		originalFileName := filepath.Base(params.File.Path)
		fileExt := filepath.Ext(originalFileName)
		fileNameWithoutExt := strings.TrimSuffix(originalFileName, fileExt)
		newFilename := fmt.Sprintf("%s.%s", fileNameWithoutExt, convertParams.Format)

		outputPath = filepath.Join(params.OutputDir, newFilename)
	} else {
		outputPath = utils.BuildOutputPath(params.OutputDir, params.File.Path)
	}

	err = params.FS.WriteFile(outputPath, newImg)
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.FinalSize, uint64(len(newImg)))
	atomic.AddUint32(stats.ProcessedImages, 1)
	return nil
}
