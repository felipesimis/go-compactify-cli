package cmd

import (
	"context"
	"fmt"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/felipesimis/go-compactify-cli/internal/validation"
	"github.com/h2non/bimg"
	"github.com/spf13/cobra"
)

type CropParams struct {
	Width   int
	Height  int
	Gravity image.Gravity
}

func cropRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	width, _ := cmd.Flags().GetInt("width")
	height, _ := cmd.Flags().GetInt("height")
	gravityInt, _ := cmd.Flags().GetInt("gravity")
	gravity := image.Gravity(gravityInt)
	dimensionValidation := &validation.DimensionsValidation{Width: width, Height: height}
	gravityValidation := &validation.GravityValidation{Gravity: gravity}
	validationComposite := validation.ValidationComposite{Validations: []validation.Validation{dimensionValidation, gravityValidation}}
	err := validationComposite.Validate()
	if err != nil {
		return err
	}
	cmd.SilenceUsage = true

	fs := filesystem.NewFileSystem()
	globalConfig := loadGlobalConfig(cmd)
	return RunOperation(globalConfig, OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
		OutputSuffix:       fmt.Sprintf("-cropped_%dx%d", width, height),
		ProgressBarMessage: "Cropping images",
		ExtraParams:        CropParams{Width: width, Height: height, Gravity: gravity},
		ProcessorFunc:      processCropImage,
		ResultVerb:         "cropped",
	})
}

func processCropImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	extraParams := params.ExtraParams.(CropParams)
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewProcessor(img)
		return newImg.Crop(extraParams.Width, extraParams.Height, extraParams.Gravity)
	})
}

var cropCmd = &cobra.Command{
	Use:     "crop",
	Aliases: []string{"cut"},
	Example: `  # Crop all images in a folder to 800x600 with center gravity
	compactify crop -i ./images -w 800 -H 600 -g 0

	# Crop and save to a specific output directory with smart gravity
	compactify crop -i ./images -o ./cropped_images -w 800 -H 600 -g 5`,
	Args:  cobra.NoArgs,
	Short: "Crop an image to specified dimensions",
	Long: `Crop an image to a specific width and height.
This command allows you to change the dimensions of an image by cropping it, which can be useful for optimizing images for 
different uses, such as web, mobile, or print. You can specify the desired width and height, 
and the image will be cropped accordingly.`,
	RunE: cropRun,
}

func init() {
	rootCmd.AddCommand(cropCmd)

	cropCmd.Flags().IntP("width", "w", 0, "Desired width of the image")
	cropCmd.Flags().IntP("height", "H", 0, "Desired height of the image")
	cropCmd.Flags().IntP("gravity", "g", int(bimg.GravityCentre), `Gravity to use when cropping the image. 
Available options:
  0 - Centre (default)
  1 - North
  2 - East
  3 - South
  4 - West
  5 - Smart`)
	cropCmd.MarkFlagRequired("width")
	cropCmd.MarkFlagRequired("height")
}
