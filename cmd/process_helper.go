package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

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
	processing.ProcessFiles(params)

	progressBar.Finish()

	totalImages := uint32(len(files))
	resultBuilder.SetTotalImages(totalImages).
		SetSkippedImages(skippedImages).
		SetProcessedImages(processedImages).
		SetOutputDirectory(outputDir).
		SetInitialSize(float64(initialSize)).
		SetFinalSize(float64(finalSize))
	result := resultBuilder.Build()
	fmt.Println(result.PrintResults(config.ResultVerb))
}
