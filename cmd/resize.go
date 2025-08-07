package cmd

import (
	"context"
	"log"

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
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewBimgImage(img)
		return newImg.Resize(extraParams.Width, extraParams.Height)
	})
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
