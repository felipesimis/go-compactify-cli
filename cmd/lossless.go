package cmd

import (
	"context"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/spf13/cobra"
)

func NewLosslessCmd(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "lossless",
		Aliases: []string{"lc"},
		Args:    cobra.NoArgs,
		Short:   "Apply lossless compression to images",
		Long: `Apply lossless compression to images.
This command allows you to apply lossless compression to images, preserving the original quality while potentially reducing the file size.
It can be useful for various image processing tasks, such as optimizing images for storage or transmission.`,
		RunE: runLossless(fs, processorFactory),
	}
}

func runLossless(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		appConfig := loadAppConfig()

		return RunOperation(appConfig, OperationConfig{
			Ctx:                ctx,
			FileSystem:         fs,
			Out:                cmd.OutOrStdout(),
			OutputSuffix:       "-lossless",
			ProgressBarMessage: "Applying lossless compression",
			ProcessorFunc: func(ctx context.Context, task processing.FileTask, stats *utils.ImageProcessingStats) error {
				return HandleImageProcessing(ctx, task, stats, processorFactory, appConfig, image.WithLosslessCompress())
			},
		})
	}
}
