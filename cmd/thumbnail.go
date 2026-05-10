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

	addEncodingFlags(cmd)
	return cmd
}

func runThumbnail(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		appConfig := loadAppConfig()

		width, _ := cmd.Flags().GetInt("width")
		validationComposite := validation.ValidationComposite{Validations: []validation.Validation{
			&validation.WidthValidation{Width: width, MinWidth: 50, MaxWidth: 1024},
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
			OutputSuffix:       "-thumbnail",
			ProgressBarMessage: "Creating thumbnails",
			ProcessorFunc: func(ctx context.Context, task processing.FileTask, stats *utils.ImageProcessingStats) error {
				return HandleImageProcessing(ctx, task, stats, processorFactory, appConfig, nil, image.WithThumbnail(width))
			},
		})
	}
}
