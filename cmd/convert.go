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

type ConvertParams struct {
	Format string
}

func convertRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	format, _ := cmd.Flags().GetString("format")

	dimensionValidation := &validation.FormatValidation{Format: format}
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
		OutputSuffix:       fmt.Sprintf("-converted.%s", format),
		ProgressBarMessage: "Converting images",
		ExtraParams:        ConvertParams{Format: format},
		ProcessorFunc:      processConvertImage,
		ResultVerb:         "converted",
	})
}

func processConvertImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	extraParams := params.ExtraParams.(ConvertParams)
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewProcessor(img)
		return newImg.Convert(extraParams.Format)
	})
}

var convertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"conv"},
	Example: `  # Convert all images in a folder to WebP
  compactify convert -i ./images -f webp

  # Convert and save to a specific output directory
  compactify convert -i ./images -o ./converted_images -f webp`,
	Args:  cobra.NoArgs,
	Short: "Convert images to a specified format",
	Long: `Convert images in a directory to a specified format.
This command allows you to change the format of images, which can be useful for optimizing images for 
different uses, such as web, mobile, or print. You can specify the desired format, 
and the images will be converted accordingly.`,
	RunE: convertRun,
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringP("format", "f", "", `Desired format of the images. Available options: webp, jpeg, png`)
	convertCmd.MarkFlagRequired("format")
}
