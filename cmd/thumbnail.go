package cmd

import (
	"context"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/image"
	"github.com/felipesimis/compactify-cli/internal/processing"
	"github.com/felipesimis/compactify-cli/internal/utils"
	"github.com/felipesimis/compactify-cli/pkg/validation"
	"github.com/spf13/cobra"
)

type ThumbnailParams struct {
	Width int
}

func thumbnailRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	dimensionValidation := &validation.WidthValidation{Width: width, MinWidth: 50, MaxWidth: 1024}
	err := dimensionValidation.Validate()
	if err != nil {
		return err
	}
	cmd.SilenceUsage = true

	fs := filesystem.NewFileSystem()
	return RunOperation(OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
		InputDir:           directory,
		OutputSuffix:       "-thumbnail",
		ProgressBarMessage: "Creating thumbnails",
		ProcessorFunc:      processThumbnailImage,
		ResultVerb:         "thumbnails created",
	})
}

func processThumbnailImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	extraParams := ThumbnailParams{Width: width}
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewBimgImage(img)
		return newImg.Thumbnail(extraParams.Width)
	})
}

var thumbnailCmd = &cobra.Command{
	Use:     "thumbnail",
	Args:    cobra.NoArgs,
	Aliases: []string{"thumb", "preview"},
	Short:   "Create a thumbnail of an image with specified width",
	Long: `Create a thumbnail of an image with a specified width, maintaining the aspect ratio 4:4.
This command allows you to generate smaller versions of images, which can be useful for previews or web usage.`,
	RunE: thumbnailRun,
}

func init() {
	rootCmd.AddCommand(thumbnailCmd)

	thumbnailCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the images to create thumbnails")
	thumbnailCmd.Flags().IntVarP(&width, "width", "w", 0, "Desired width of the thumbnail")

	thumbnailCmd.MarkFlagRequired("directory")
	thumbnailCmd.MarkFlagRequired("width")
}
