package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/image"
	"github.com/felipesimis/compactify-cli/internal/processing"
	"github.com/felipesimis/compactify-cli/internal/utils"
	"github.com/felipesimis/compactify-cli/pkg/progress"
	"github.com/spf13/cobra"
)

func flipRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fs := filesystem.NewFileSystem()
	files, err := fs.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	outputDir, err := fs.CreateSiblingDir(directory, "-flipped")
	if err != nil {
		log.Fatal(err)
	}

	var initialSize, finalSize uint64
	var skippedImages, flippedImages uint32

	resultBuilder := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(os.Stdout, len(files), concurrency, "Resizing images")

	params := processing.ProcessFilesParams{
		Files:       files,
		FS:          fs,
		OutputDir:   outputDir,
		ProgressBar: progressBar,
		ProcessorFunc: func(p processing.FileProcessingParams) error {
			stats := utils.NewImageProcessingStats(&initialSize, &finalSize, &skippedImages, &flippedImages)
			return flipImages(ctx, p, stats)
		},
		Concurrency: concurrency,
	}
	processing.ProcessFiles(params)

	progressBar.Finish()

	totalImages := uint32(len(files))
	resultBuilder.SetTotalImages(totalImages).
		SetSkippedImages(skippedImages).
		SetProcessedImages(flippedImages).
		SetOutputDirectory(outputDir).
		SetInitialSize(float64(initialSize)).
		SetFinalSize(float64(finalSize))
	result := resultBuilder.Build()
	fmt.Println(result.PrintResults("flipped"))
}

func flipImages(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	img, err := params.FS.ReadFile(params.File.Path)
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.InitialSize, uint64(params.File.Size))
	newImg := image.NewBimgImage(img)
	flippedImg, err := newImg.Flip()
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	outputPath := utils.BuildOutputPath(params.OutputDir, params.File.Path)
	err = params.FS.WriteFile(outputPath, flippedImg)
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.FinalSize, uint64(len(flippedImg)))
	atomic.AddUint32(stats.ProcessedImages, 1)
	return nil
}

var flipCmd = &cobra.Command{
	Use:     "flip",
	Aliases: []string{"invert", "mirror"},
	Args:    cobra.NoArgs,
	Short:   "Flip images vertically",
	Long: `Flip images vertically.
This command allows you to flip an image along the vertical axis, creating a mirror image.
It can be useful for various image processing tasks, such as creating reflections or correcting image orientation.`,
	Run: flipRun,
}

func init() {
	rootCmd.AddCommand(flipCmd)

	flipCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the images to flip")
	flipCmd.MarkFlagRequired("directory")
}
