package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

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

func (c ConvertParams) ModifyOutputPath(relPath, outputDir string) string {
	if c.Format == "" {
		return ""
	}

	fileExt := filepath.Ext(relPath)
	relPathWithoutExt := strings.TrimSuffix(relPath, fileExt)
	newFilename := fmt.Sprintf("%s.%s", relPathWithoutExt, c.Format)
	return filepath.Join(outputDir, newFilename)
}

func NewConvertCmd(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) *cobra.Command {
	cmd := &cobra.Command{
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
		RunE: runConvert(fs, processorFactory),
	}

	cmd.Flags().StringP("format", "f", "", `Desired format of the images. Available options: webp, jpeg, png`)
	cmd.MarkFlagRequired("format")
	addEncodingFlags(cmd)

	return cmd
}

func runConvert(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		appConfig := loadAppConfig()

		format, _ := cmd.Flags().GetString("format")
		validationComposite := validation.ValidationComposite{Validations: []validation.Validation{
			&validation.FormatValidation{Format: format},
			&validation.QualityValidation{Quality: appConfig.Quality},
		}}
		if err := validationComposite.Validate(); err != nil {
			return err
		}
		cmd.SilenceUsage = true

		return RunOperation(appConfig, OperationConfig{
			Ctx:                ctx,
			FileSystem:         fs,
			Out:                cmd.OutOrStdout(),
			OutputSuffix:       fmt.Sprintf("-converted.%s", format),
			ProgressBarMessage: "Converting images",
			ExtraParams:        ConvertParams{Format: format},
			ProcessorFunc: func(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
				extraParams := params.ExtraParams.(ConvertParams)
				return HandleImageProcessing(ctx, params, stats, processorFactory, appConfig, image.WithConvert(extraParams.Format))
			},
		})
	}
}
