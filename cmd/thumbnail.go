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
	"github.com/felipesimis/compactify-cli/pkg/validation"
	"github.com/spf13/cobra"
)

type ThumbnailParams struct {
	Width int
}

type ThumbnailStats struct {
	initialSize     *uint64
	finalSize       *uint64
	skippedImages   *uint32
	thumbnailImages *uint32
}

func thumbnailRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dimensionValidation := &validation.WidthValidation{Width: width, MinWidth: 50, MaxWidth: 1024}
	err := dimensionValidation.Validate()
	if err != nil {
		log.Fatal(err)
	}

	fs := filesystem.NewFileSystem()
	files, err := fs.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	outputDir, err := fs.CreateSiblingDir(directory, "-thumbnails")
	if err != nil {
		log.Fatal(err)
	}

	var initialSize, finalSize uint64
	var skippedImages, thumbnailImages uint32

	resultBuilder := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(os.Stdout, len(files), concurrency, "Creating thumbnails")

	params := processing.ProcessFilesParams{
		Files:       files,
		FS:          fs,
		OutputDir:   outputDir,
		ProgressBar: progressBar,
		ExtraParams: ThumbnailParams{Width: width},
		ProcessorFunc: func(p processing.FileProcessingParams) error {
			extraParams := p.ExtraParams.(ThumbnailParams)
			stats := &ThumbnailStats{
				initialSize:     &initialSize,
				finalSize:       &finalSize,
				skippedImages:   &skippedImages,
				thumbnailImages: &thumbnailImages,
			}
			return processThumbnailImages(ctx, p, extraParams, stats)
		},
		Concurrency: concurrency,
	}
	processing.ProcessFiles(params)

	progressBar.Finish()

	totalImages := uint32(len(files))
	resultBuilder.SetTotalImages(totalImages).
		SetSkippedImages(skippedImages).
		SetProcessedImages(thumbnailImages).
		SetOutputDirectory(outputDir).
		SetInitialSize(float64(initialSize)).
		SetFinalSize(float64(finalSize))
	result := resultBuilder.Build()
	fmt.Println(result.PrintResults("thumbnails"))
}

func processThumbnailImages(ctx context.Context, params processing.FileProcessingParams, extraParams ThumbnailParams, stats *ThumbnailStats) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	img, err := params.FS.ReadFile(params.File.Path)
	if err != nil {
		atomic.AddUint32(stats.skippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.initialSize, uint64(params.File.Size))
	newImg := image.NewBimgImage(img)
	thumbnailsImg, err := newImg.Thumbnail(extraParams.Width)
	if err != nil {
		atomic.AddUint32(stats.skippedImages, 1)
		return err
	}

	outputPath := utils.BuildOutputPath(params.OutputDir, params.File.Path)
	err = params.FS.WriteFile(outputPath, thumbnailsImg)
	if err != nil {
		atomic.AddUint32(stats.skippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.finalSize, uint64(len(thumbnailsImg)))
	atomic.AddUint32(stats.thumbnailImages, 1)
	return nil
}

var thumbnailCmd = &cobra.Command{
	Use:   "thumbnail",
	Args:  cobra.NoArgs,
	Short: "Create a thumbnail of an image with specified width",
	Long: `Create a thumbnail of an image with a specified width, maintaining the aspect ratio 4:4.
This command allows you to generate smaller versions of images, which can be useful for previews or web usage.`,
	Run: thumbnailRun,
}

func init() {
	rootCmd.AddCommand(thumbnailCmd)

	thumbnailCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the images to create thumbnails")
	thumbnailCmd.Flags().IntVarP(&width, "width", "w", 0, "Desired width of the thumbnail")

	thumbnailCmd.MarkFlagRequired("directory")
	thumbnailCmd.MarkFlagRequired("width")
}
