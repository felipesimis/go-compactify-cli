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
		buf := bytes.NewBuffer(make([]byte, 0, 5*bytesInMb))
		return buf
	},
}

type OutputPathModifier interface {
	ModifyOutputPath(relPath, outputDir string) string
}

type OperationConfig struct {
	Ctx                context.Context
	FileSystem         filesystem.FileSystem
	Out                io.Writer
	OutputSuffix       string
	ProgressBarMessage string
	ProcessorFunc      func(ctx context.Context, task processing.FileTask, stats *utils.ImageProcessingStats) error
}

func RunOperation(app AppConfig, config OperationConfig) error {
	if app.DryRun {
		config.FileSystem = filesystem.NewDryRunFileSystem(config.FileSystem)
		fmt.Fprintln(config.Out, ui.Warn("DRY-RUN MODE: No files will be modified or created on disk."))
	}

	outputRootDir, err := resolveRootOutputDir(app, config)
	if err != nil {
		return err
	}

	if app.Recursive {
		fmt.Fprintln(config.Out, ui.Info(fmt.Sprintf("Analyzing directories recursively in %s ...", app.InputDir)))
	} else {
		fmt.Fprintln(config.Out, ui.Info(fmt.Sprintf("Analyzing files in %s ...", app.InputDir)))
	}

	files, err := processing.DiscoverAndPrepare(config.FileSystem, app.InputDir, outputRootDir, app.Recursive)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Fprintln(config.Out, ui.Warn(fmt.Sprintf("No files found in directory: %s", app.InputDir)))
		return nil
	}

	stats := &utils.ImageProcessingStats{}
	resultBuilder := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(config.Out, len(files), app.Concurrency, config.ProgressBarMessage)

	wrappedProcessor := func(task processing.FileTask) error {
		return config.ProcessorFunc(config.Ctx, task, stats)
	}
	fileBatchConfig := processing.FileBatchConfig{
		Files:       files,
		FS:          config.FileSystem,
		InputDir:    app.InputDir,
		OutputDir:   outputRootDir,
		ProgressBar: progressBar,
		Handler:     wrappedProcessor,
		Concurrency: app.Concurrency,
	}
	processErrors := processing.RunFileBatch(fileBatchConfig)
	progressBar.Finish()
	totalImages := uint32(len(files))
	resultBuilder.SetTotalImages(totalImages).
		SetSkippedImages(stats.SkippedImages.Load()).
		SetProcessedImages(stats.ProcessedImages.Load()).
		SetOutputDirectory(outputRootDir).
		SetOriginalBytes(stats.InitialSize.Load()).
		SetProcessedBytes(stats.FinalSize.Load()).
		SetErrors(processErrors)
	result := resultBuilder.Build()
	fmt.Fprintln(config.Out, RenderProcessSummary(result))

	return nil
}

func HandleImageProcessing(
	ctx context.Context,
	task processing.FileTask,
	stats *utils.ImageProcessingStats,
	processorFactory image.ProcessorFactory,
	appConfig AppConfig,
	modifier OutputPathModifier,
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
	buf.Grow(int(task.File.Size))
	defer bufferPool.Put(buf)

	file, err := task.FS.OpenFile(task.File.Path)
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
	finalOpts = append(finalOpts, buildAppOptions(appConfig)...)

	processor := processorFactory(imgBytes)
	newImg, err := processor.Process(finalOpts...)
	if err != nil {
		stats.SkippedImages.Add(1)
		return err
	}

	outputPath := resolveImageDestination(task.File, task.OutputDir, modifier)

	destDir := filepath.Dir(outputPath)
	if err := task.FS.CreateDir(destDir); err != nil {
		stats.SkippedImages.Add(1)
		return err
	}

	err = task.FS.WriteFile(outputPath, newImg)
	if err != nil {
		stats.SkippedImages.Add(1)
		return err
	}

	stats.FinalSize.Add(uint64(len(newImg)))
	stats.ProcessedImages.Add(1)
	return nil
}

func buildAppOptions(appConfig AppConfig) []image.ProcessOption {
	var opts []image.ProcessOption
	if appConfig.StripMetadata {
		opts = append(opts, image.WithStripMetadata())
	}
	if appConfig.Quality > 0 {
		opts = append(opts, image.WithQuality(appConfig.Quality))
	}
	return opts
}

func resolveRootOutputDir(appConfig AppConfig, config OperationConfig) (string, error) {
	if appConfig.OutputDir != "" {
		if err := config.FileSystem.CreateDir(appConfig.OutputDir); err != nil {
			return "", err
		}
		return appConfig.OutputDir, nil
	}
	return config.FileSystem.CreateSiblingDir(appConfig.InputDir, config.OutputSuffix)
}

func resolveImageDestination(file filesystem.FileInfo, outputDir string, modifier OutputPathModifier) string {
	relPath := file.RelPath
	if relPath == "" {
		relPath = filepath.Base(file.Path)
	}

	if modifier != nil {
		return modifier.ModifyOutputPath(relPath, outputDir)
	}

	return utils.BuildOutputPath(outputDir, relPath)
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
