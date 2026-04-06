package cmd

import (
	"context"
	"fmt"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/image"
	"github.com/felipesimis/compactify-cli/internal/processing"
	"github.com/felipesimis/compactify-cli/internal/utils"
	"github.com/felipesimis/compactify-cli/pkg/validation"
	"github.com/h2non/bimg"
	"github.com/spf13/cobra"
)

var (
	gravity int
)

type CropParams struct {
	Width   int
	Height  int
	Gravity bimg.Gravity
}

func cropRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	dimensionValidation := &validation.DimensionsValidation{Width: width, Height: height}
	gravityValidation := &validation.GravityValidation{Gravity: bimg.Gravity(gravity)}
	validationComposite := validation.ValidationComposite{Validations: []validation.Validation{dimensionValidation, gravityValidation}}
	err := validationComposite.Validate()
	if err != nil {
		return err
	}
	cmd.SilenceUsage = true

	fs := filesystem.NewFileSystem()

	return RunOperation(OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
		InputDir:           directory,
		OutputSuffix:       fmt.Sprintf("-cropped_%dx%d", width, height),
		ProgressBarMessage: "Cropping images",
		ExtraParams:        CropParams{Width: width, Height: height, Gravity: bimg.Gravity(gravity)},
		ProcessorFunc:      processCropImage,
		ResultVerb:         "cropped",
	})
}

func processCropImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	extraParams := params.ExtraParams.(CropParams)
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewBimgImage(img)
		return newImg.Crop(extraParams.Width, extraParams.Height, extraParams.Gravity)
	})
}

var cropCmd = &cobra.Command{
	Use:     "crop",
	Aliases: []string{"cut"},
	Args:    cobra.NoArgs,
	Short:   "Crop an image to specified dimensions",
	Long: `Crop an image to a specific width and height.
This command allows you to change the dimensions of an image by cropping it, which can be useful for optimizing images for 
different uses, such as web, mobile, or print. You can specify the desired width and height, 
and the image will be cropped accordingly.`,
	RunE: cropRun,
}

func init() {
	rootCmd.AddCommand(cropCmd)

	cropCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the image to crop")
	cropCmd.Flags().IntVarP(&width, "width", "w", 0, "Desired width of the image")
	cropCmd.Flags().IntVarP(&height, "height", "H", 0, "Desired height of the image")
	cropCmd.Flags().IntVarP(&gravity, "gravity", "g", int(bimg.GravityCentre), `Gravity to use when cropping the image. 
Available options:
  0 - Centre (default)
  1 - North
  2 - East
  3 - South
  4 - West
  5 - Smart`)

	cropCmd.MarkFlagRequired("directory")
	cropCmd.MarkFlagRequired("width")
	cropCmd.MarkFlagRequired("height")
}
