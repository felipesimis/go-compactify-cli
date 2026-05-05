package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sync"
	"time"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/ui"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/felipesimis/go-compactify-cli/pkg/progress"
)

const (
	bytesInMb = 1024 * 1024
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type OutputPathModifier interface {
	ModifyOutputPath(originalPath, outputDir string) string
}

type OperationConfig struct {
	Ctx                context.Context
	FileSystem         filesystem.FileSystem
	Out                io.Writer
	OutputSuffix       string
	ProgressBarMessage string
	ExtraParams        interface{}
	ProcessorFunc      func(ctx context.Context, p processing.FileProcessingParams, stats *utils.ImageProcessingStats) error
}

func RunOperation(global GlobalConfig, config OperationConfig) error {
	if global.DryRun {
		config.FileSystem = filesystem.NewDryRunFileSystem(config.FileSystem)

		fmt.Fprintln(config.Out, ui.Warn("DRY-RUN MODE: No files will be modified or created on disk."))
	}

	files, err := config.FileSystem.ReadDir(global.InputDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Fprintln(config.Out, ui.Warn(fmt.Sprintf("No files found in directory: %s", global.InputDir)))
		return nil
	}

	finalOutputDir, err := resolveOutputDir(global, config)
	if err != nil {
		return err
	}

	stats := &utils.ImageProcessingStats{}
	resultBuilder := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(config.Out, len(files), global.Concurrency, config.ProgressBarMessage)
	defer progressBar.Finish()

	wrappedProcessor := func(p processing.FileProcessingParams) error {
		return config.ProcessorFunc(config.Ctx, p, stats)
	}
	params := processing.ProcessFilesParams{
		Files:         files,
		FS:            config.FileSystem,
		InputDir:      global.InputDir,
		OutputDir:     finalOutputDir,
		ProgressBar:   progressBar,
		ExtraParams:   config.ExtraParams,
		ProcessorFunc: wrappedProcessor,
		Concurrency:   global.Concurrency,
	}
	processErrors := processing.ProcessFiles(params)
	totalImages := uint32(len(files))
	resultBuilder.SetTotalImages(totalImages).
		SetSkippedImages(stats.SkippedImages.Load()).
		SetProcessedImages(stats.ProcessedImages.Load()).
		SetOutputDirectory(finalOutputDir).
		SetOriginalBytes(stats.InitialSize.Load()).
		SetProcessedBytes(stats.FinalSize.Load()).
		SetErrors(processErrors)
	result := resultBuilder.Build()
	fmt.Fprintln(config.Out, RenderProcessSummary(result))

	return nil
}

func HandleImageProcessing(
	ctx context.Context,
	params processing.FileProcessingParams,
	stats *utils.ImageProcessingStats,
	processorFactory image.ProcessorFactory,
	globalCfg GlobalConfig,
	opts ...image.ProcessOption,
) error {
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

	finalOpts := append([]image.ProcessOption{}, opts...)
	finalOpts = append(finalOpts, buildGlobalOptions(globalCfg)...)

	processor := processorFactory(imgBytes)
	newImg, err := processor.Process(finalOpts...)
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

func buildGlobalOptions(global GlobalConfig) []image.ProcessOption {
	var opts []image.ProcessOption
	if global.StripMetadata {
		opts = append(opts, image.WithStripMetadata())
	}
	return opts
}

func resolveOutputDir(global GlobalConfig, config OperationConfig) (string, error) {
	if global.OutputDir != "" {
		if err := config.FileSystem.CreateDir(global.OutputDir); err != nil {
			return "", err
		}
		return global.OutputDir, nil
	}
	return config.FileSystem.CreateSiblingDir(global.InputDir, config.OutputSuffix)
}

func determineOutputPath(params processing.FileProcessingParams) string {
	if modifier, ok := params.ExtraParams.(OutputPathModifier); ok {
		return modifier.ModifyOutputPath(params.File.Path, params.OutputDir)
	}

	relativePath, err := filepath.Rel(params.InputDir, params.File.Path)
	if err != nil {
		relativePath = filepath.Base(params.File.Path)
	}

	return utils.BuildOutputPath(params.OutputDir, relativePath)
}

func RenderProcessSummary(r *utils.Result) string {
	items := []ui.Item{
		{Label: "Time", Value: r.ElapsedTime.Round(time.Millisecond).String()},
		{Label: "Total", Value: fmt.Sprintf("%d images", r.TotalImages)},
	}

	if r.SkippedImages > 0 {
		items = append(items, ui.Item{
			Label: "Skipped",
			Value: fmt.Sprintf("%d", r.SkippedImages),
		})
	}

	left := ui.Panel{
		Title: "OPERATION",
		Items: append(items, ui.Item{
			Label: "Processed",
			Value: fmt.Sprintf("%d", r.ProcessedImages),
		}),
	}

	toMB := func(b int64) float64 { return float64(b) / bytesInMb }
	right := ui.Panel{
		Title: "IMPACT",
		Items: []ui.Item{
			{Label: "Original", Value: fmt.Sprintf("%.2f MB", toMB(r.OriginalBytes)), IsHighlighted: false},
			{Label: "After", Value: fmt.Sprintf("%.2f MB", toMB(r.ProcessedBytes)), IsHighlighted: false},
			{Label: "", Value: ""},
			{Label: "Saved", Value: fmt.Sprintf("%.2f MB", toMB(r.SavedBytes)), IsHighlighted: true},
			{Label: "Reduction", Value: fmt.Sprintf("%.2f%%", r.ReductionRatio), IsHighlighted: true},
		},
	}

	dashboard := ui.RenderDashboard(left, right, "OUTPUT DIRECTORY", fmt.Sprintf("📂 %s", r.OutputDirectory))
	errors := ui.RenderErrorList(r.Errors)

	return "\n" + dashboard + errors + "\n"
}
