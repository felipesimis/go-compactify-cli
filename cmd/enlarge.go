package cmd

import (
	"context"
	"fmt"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/felipesimis/go-compactify-cli/internal/validation"
	"github.com/spf13/cobra"
)

type EnlargeParams struct {
	Width  int
	Height int
}

func enlargeRun(cmd *cobra.Command, args []string) error {
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
		OutputSuffix:       fmt.Sprintf("-enlarged-%dx%d", width, height),
		ProgressBarMessage: "Enlarging images",
		ExtraParams:        EnlargeParams{Width: width, Height: height},
		ProcessorFunc:      processEnlargeImage,
		ResultVerb:         "enlarged",
	})
}

func processEnlargeImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	extraParams := params.ExtraParams.(EnlargeParams)
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewProcessor(img)
		return newImg.Enlarge(extraParams.Width, extraParams.Height)
	})
}

var enlargeCmd = &cobra.Command{
	Use: "enlarge",
	Example: `  # Enlarge all images in a folder to 1200x900
  compactify enlarge -i ./images -w 1200 -H 900

  # Enlarge and save to a specific output directory
  compactify enlarge -i ./images -o ./enlarged_images -w 1200 -H 900`,
	Args:  cobra.NoArgs,
	Short: "Enlarge an image to specified dimensions while maintaining aspect ratio",
	Long: `Enlarge an image to a specific width and height while maintaining the aspect ratio.
This command allows you to change the dimensions of an image, which can be useful for optimizing images for 
different uses, such as web, mobile, or print. You can specify the desired width and height, 
and the image will be enlarged accordingly, keeping its original aspect ratio.`,
	RunE: enlargeRun,
}

func init() {
	rootCmd.AddCommand(enlargeCmd)

	enlargeCmd.Flags().IntP("width", "w", 0, "Desired width of the image")
	enlargeCmd.Flags().IntP("height", "H", 0, "Desired height of the image")
	enlargeCmd.MarkFlagRequired("width")
	enlargeCmd.MarkFlagRequired("height")
}
