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

var (
	directory string
	width     int
	height    int
)

type ResizeParams struct {
	Width  int
	Height int
}

type ResizeStats struct {
	initialSize   *uint64
	finalSize     *uint64
	skippedImages *uint32
	resizedImages *uint32
}

func resizeRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dimensionValidation := &validation.DimensionsValidation{Width: width, Height: height}
	err := dimensionValidation.Validate()
	if err != nil {
		log.Fatal(err)
	}

	fs := filesystem.NewFileSystem()
	files, err := fs.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	outputDir, err := fs.CreateSiblingDir(directory, "-resized")
	if err != nil {
		log.Fatal(err)
	}

	var initialSize, finalSize uint64
	var skippedImages, resizedImages uint32

	resultBuilder := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(os.Stdout, len(files), concurrency, "Resizing images")

	params := processing.ProcessFilesParams{
		Files:       files,
		FS:          fs,
		OutputDir:   outputDir,
		ProgressBar: progressBar,
		ExtraParams: ResizeParams{Width: width, Height: height},
		ProcessorFunc: func(p processing.FileProcessingParams) error {
			extraParams := p.ExtraParams.(ResizeParams)
			stats := &ResizeStats{
				initialSize:   &initialSize,
				finalSize:     &finalSize,
				skippedImages: &skippedImages,
				resizedImages: &resizedImages,
			}
			return resizeImages(ctx, p, extraParams, stats)
		},
		Concurrency: concurrency,
	}
	processing.ProcessFiles(params)

	progressBar.Finish()

	totalImages := uint32(len(files))
	resultBuilder.SetTotalImages(totalImages).
		SetSkippedImages(skippedImages).
		SetProcessedImages(resizedImages).
		SetOutputDirectory(outputDir).
		SetInitialSize(float64(initialSize)).
		SetFinalSize(float64(finalSize))
	result := resultBuilder.Build()
	fmt.Println(result.PrintResults("resized"))
}

func resizeImages(ctx context.Context, params processing.FileProcessingParams, extraParams ResizeParams, stats *ResizeStats) error {
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
	resizedImg, err := newImg.Resize(extraParams.Width, extraParams.Height)
	if err != nil {
		atomic.AddUint32(stats.skippedImages, 1)
		return err
	}

	outputPath := utils.BuildOutputPath(params.OutputDir, params.File.Path)
	err = params.FS.WriteFile(outputPath, resizedImg)
	if err != nil {
		atomic.AddUint32(stats.skippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.finalSize, uint64(len(resizedImg)))
	atomic.AddUint32(stats.resizedImages, 1)
	return nil
}

var resizeCmd = &cobra.Command{
	Use:     "resize",
	Aliases: []string{"scale", "rescale"},
	Args:    cobra.NoArgs,
	Short:   "Resize an image to specified dimensions",
	Long: `Resize an image to a specific width and height.
This command allows you to change the dimensions of an image, which can be useful for optimizing images for 
different uses, such as web, mobile, or print. You can specify the desired width and height, 
and the image will be resized accordingly.`,
	Run: resizeRun,
}

func init() {
	rootCmd.AddCommand(resizeCmd)

	resizeCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the images to resize")
	resizeCmd.Flags().IntVarP(&width, "width", "w", 0, "Desired width of the image")
	resizeCmd.Flags().IntVarP(&height, "height", "H", 0, "Desired height of the image")

	resizeCmd.MarkFlagRequired("directory")
	resizeCmd.MarkFlagRequired("width")
	resizeCmd.MarkFlagRequired("height")
}
