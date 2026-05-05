package cmd

import (
	"context"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/spf13/cobra"
)

func NewGrayscaleCmd(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "grayscale",
		Aliases: []string{"gray", "bw"},
		Args:    cobra.NoArgs,
		Short:   "Convert images to grayscale",
		Long: `Convert images to grayscale.
This command allows you to convert an image to grayscale, removing all color information and leaving only shades of gray.
It can be useful for various image processing tasks, such as creating artistic effects or preparing images for printing.`,
		RunE: runGrayscale(fs, processorFactory),
	}
}

func runGrayscale(fs filesystem.FileSystem, processorFactory image.ProcessorFactory) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		globalConfig := loadGlobalConfig(cmd)

		return RunOperation(globalConfig, OperationConfig{
			Ctx:                ctx,
			FileSystem:         fs,
			Out:                cmd.OutOrStdout(),
			OutputSuffix:       "-grayscale",
			ProgressBarMessage: "Creating grayscale images",
			ProcessorFunc: func(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
				return HandleImageProcessing(ctx, params, stats, processorFactory, globalConfig, image.WithGrayscale())
			},
		})
	}
}
