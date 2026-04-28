package cmd

import (
	"context"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/spf13/cobra"
)

func losslessRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	fs := filesystem.NewFileSystem()
	globalConfig := loadGlobalConfig(cmd)
	return RunOperation(globalConfig, OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
		OutputSuffix:       "-lossless",
		ProgressBarMessage: "Applying lossless compression",
		ProcessorFunc:      processLosslessImage,
		ResultVerb:         "lossless compressed",
	})
}

func processLosslessImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewProcessor(img)
		return newImg.LosslessCompress()
	})
}

var losslessCmd = &cobra.Command{
	Use:     "lossless",
	Aliases: []string{"lc"},
	Args:    cobra.NoArgs,
	Short:   "Apply lossless compression to images",
	Long: `Apply lossless compression to images.
This command allows you to apply lossless compression to images, preserving the original quality while potentially reducing the file size.
It can be useful for various image processing tasks, such as optimizing images for storage or transmission.`,
	RunE: losslessRun,
}

func init() {
	rootCmd.AddCommand(losslessCmd)
}
