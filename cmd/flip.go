package cmd

import (
	"context"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/spf13/cobra"
)

func flipRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	fs := filesystem.NewFileSystem()
	globalConfig := loadGlobalConfig(cmd)
	return RunOperation(globalConfig, OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
		OutputSuffix:       "-flipped",
		ProgressBarMessage: "Flipping images",
		ProcessorFunc:      processFlipImage,
		ResultVerb:         "flipped",
	})
}

func processFlipImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewBimgImage(img)
		return newImg.Flip()
	})
}

var flipCmd = &cobra.Command{
	Use:     "flip",
	Aliases: []string{"invert", "mirror"},
	Args:    cobra.NoArgs,
	Short:   "Flip images vertically",
	Long: `Flip images vertically.
This command allows you to flip an image along the vertical axis, creating a mirror image.
It can be useful for various image processing tasks, such as creating reflections or correcting image orientation.`,
	RunE: flipRun,
}

func init() {
	rootCmd.AddCommand(flipCmd)
}
