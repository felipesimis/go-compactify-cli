package cmd

import (
	"context"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/felipesimis/go-compactify-cli/internal/validation"
	"github.com/spf13/cobra"
)

type ResizeParams struct {
	Width  int
	Height int
}

func resizeRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	width, _ := cmd.Flags().GetInt("width")
	height, _ := cmd.Flags().GetInt("height")

	dimensionValidation := &validation.DimensionsValidation{Width: width, Height: height}
	err := dimensionValidation.Validate()
	if err != nil {
		return err
	}
	cmd.SilenceUsage = true

	fs := filesystem.NewFileSystem()
	globalConfig := loadGlobalConfig(cmd)
	return RunOperation(globalConfig, OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
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
	Example: `  # Resize all images in a folder to 800x600
  compactify resize -i ./images -w 800 -H 600

  # Resize and save to a specific output directory
  compactify resize -i ./images -o ./resized_images -w 800 -H 600`,
	Args:  cobra.NoArgs,
	Short: "Resize an image to specified dimensions",
	Long: `Resize an image to a specific width and height.
This command allows you to change the dimensions of an image, which can be useful for optimizing images for 
different uses, such as web, mobile, or print. You can specify the desired width and height, 
and the image will be resized accordingly.`,
	RunE: resizeRun,
}

func init() {
	rootCmd.AddCommand(resizeCmd)

	resizeCmd.Flags().IntP("width", "w", 0, "Desired width of the image")
	resizeCmd.Flags().IntP("height", "H", 0, "Desired height of the image")
	resizeCmd.MarkFlagRequired("width")
	resizeCmd.MarkFlagRequired("height")
}
