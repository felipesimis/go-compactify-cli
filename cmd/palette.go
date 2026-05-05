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

func NewPaletteCmd(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "palette",
		Args:  cobra.NoArgs,
		Short: "Enable palette on images",
		Long: `Apply a color palette to images.
This command enables a color palette on the specified images, which can help reduce the file size by limiting the number of colors used. 
It is useful for optimizing images for web use, creating artistic effects, and ensuring compatibility with formats that require or benefit from a limited color palette.`,
		RunE: runPalette(fs, processorFactory),
	}

	addEncodingFlags(cmd)
	return cmd
}

func runPalette(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		appConfig := loadAppConfig()

		qualityValidation := &validation.QualityValidation{Quality: appConfig.Quality}
		if err := qualityValidation.Validate(); err != nil {
			return err
		}

		return RunOperation(appConfig, OperationConfig{
			Ctx:                ctx,
			FileSystem:         fs,
			Out:                cmd.OutOrStdout(),
			OutputSuffix:       "-palette",
			ProgressBarMessage: "Enabling palette on images",
			ProcessorFunc: func(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
				return HandleImageProcessing(ctx, params, stats, processorFactory, appConfig,
					image.WithPalette(),
					image.WithQuality(appConfig.Quality),
				)
			},
		})
	}
}
