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

func paletteRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fs := filesystem.NewFileSystem()
	files, err := fs.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	outputDir, err := fs.CreateSiblingDir(directory, "-palette")
	if err != nil {
		log.Fatal(err)
	}

	var initialSize, finalSize uint64
	var skippedImages, processedImages uint32

	resultBuilder := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(os.Stdout, len(files), concurrency, "Applying palette on images")

	params := processing.ProcessFilesParams{
		Files:       files,
		FS:          fs,
		OutputDir:   outputDir,
		ProgressBar: progressBar,
		ProcessorFunc: func(p processing.FileProcessingParams) error {
			stats := utils.NewImageProcessingStats(&initialSize, &finalSize, &skippedImages, &processedImages)
			return processPaletteImage(ctx, p, stats)
		},
		Concurrency: concurrency,
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
	fmt.Println(result.PrintResults("palette"))
}

func processPaletteImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
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
	paletteImageData, err := newImg.EnablePalette()
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	outputPath := utils.BuildOutputPath(params.OutputDir, params.File.Path)
	err = params.FS.WriteFile(outputPath, paletteImageData)
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.FinalSize, uint64(len(paletteImageData)))
	atomic.AddUint32(stats.ProcessedImages, 1)
	return nil
}

var paletteCmd = &cobra.Command{
	Use:   "palette",
	Args:  cobra.NoArgs,
	Short: "Enable palette on images",
	Long: `Apply a color palette to images.
This command enables a color palette on the specified images, which can help reduce the file size by limiting the number of colors used. 
It is useful for optimizing images for web use, creating artistic effects, and ensuring compatibility with formats that require or benefit from a limited color palette.`,
	Run: paletteRun,
}

func init() {
	rootCmd.AddCommand(paletteCmd)

	paletteCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the images to apply palette")
	paletteCmd.MarkFlagRequired("directory")
}
