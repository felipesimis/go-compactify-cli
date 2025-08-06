package cmd

import (
	"context"
	"log"
	"sync/atomic"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/image"
	"github.com/felipesimis/compactify-cli/internal/processing"
	"github.com/felipesimis/compactify-cli/internal/utils"
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

func resizeRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dimensionValidation := &validation.DimensionsValidation{Width: width, Height: height}
	err := dimensionValidation.Validate()
	if err != nil {
		log.Fatal(err)
	}

	fs := filesystem.NewFileSystem()

	RunOperation(OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
		InputDir:           directory,
		OutputSuffix:       "-resized",
		ProgressBarMessage: "Resizing images",
		ExtraParams:        ResizeParams{Width: width, Height: height},
		ProcessorFunc:      processResizeImage,
		ResultVerb:         "resized",
	})
}

func processResizeImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	extraParams := params.ExtraParams.(ResizeParams)

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
	newImg := image.NewBimgImage(img)
	resizedImg, err := newImg.Resize(extraParams.Width, extraParams.Height)
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	outputPath := utils.BuildOutputPath(params.OutputDir, params.File.Path)
	err = params.FS.WriteFile(outputPath, resizedImg)
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.FinalSize, uint64(len(resizedImg)))
	atomic.AddUint32(stats.ProcessedImages, 1)
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
