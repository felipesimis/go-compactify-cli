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

func NewFlipCmd(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "flip",
		Aliases: []string{"invert", "mirror"},
		Args:    cobra.NoArgs,
		Short:   "Flip images vertically",
		Long: `Flip images vertically.
This command allows you to flip an image along the vertical axis, creating a mirror image.
It can be useful for various image processing tasks, such as creating reflections or correcting image orientation.`,
		RunE: runFlip(fs, processorFactory),
	}

	addEncodingFlags(cmd)
	return cmd
}

func runFlip(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) func(cmd *cobra.Command, args []string) error {
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
			OutputSuffix:       "-flipped",
			ProgressBarMessage: "Flipping images",
			ProcessorFunc: func(ctx context.Context, task processing.FileTask, stats *utils.ImageProcessingStats) error {
				return HandleImageProcessing(ctx, task, stats, processorFactory, appConfig, nil, image.WithFlip())
			},
		})
	}
}
