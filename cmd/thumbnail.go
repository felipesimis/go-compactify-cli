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

type ThumbnailParams struct {
	Width int
}

func NewThumbnailCmd(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "thumbnail",
		Args:    cobra.NoArgs,
		Aliases: []string{"thumb", "preview"},
		Example: `  # Create thumbnails for all images in a folder
  compactify thumbnail -i ./images -w 150

	# Create thumbnails and save to a specific output directory
	compactify thumbnail -i ./images -o ./thumbnails -w 150`,
		Short: "Create a thumbnail of an image with specified width",
		Long: `Create a thumbnail of an image with a specified width, maintaining the aspect ratio 4:4.
This command allows you to generate smaller versions of images, which can be useful for previews or web usage.`,
		RunE: runThumbnail(fs, processorFactory),
	}

	cmd.Flags().IntP("width", "w", 0, "Desired width of the thumbnail")
	cmd.MarkFlagRequired("width")

	return cmd
}

func runThumbnail(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		width, _ := cmd.Flags().GetInt("width")
		dimensionValidation := &validation.WidthValidation{Width: width, MinWidth: 50, MaxWidth: 1024}
		err := dimensionValidation.Validate()
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		appConfig := loadAppConfig()

		return RunOperation(appConfig, OperationConfig{
			Ctx:                ctx,
			FileSystem:         fs,
			Out:                cmd.OutOrStdout(),
			OutputSuffix:       "-thumbnail",
			ProgressBarMessage: "Creating thumbnails",
			ExtraParams:        ThumbnailParams{Width: width},
			ProcessorFunc: func(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
				extraParams := params.ExtraParams.(ThumbnailParams)
				return HandleImageProcessing(ctx, params, stats, processorFactory, appConfig, image.WithThumbnail(extraParams.Width))
			},
		})
	}
}
