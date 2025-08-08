package cmd

import (
	"context"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/image"
	"github.com/felipesimis/compactify-cli/internal/processing"
	"github.com/felipesimis/compactify-cli/internal/utils"
	"github.com/spf13/cobra"
)

func flipRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fs := filesystem.NewFileSystem()
	RunOperation(OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
		InputDir:           directory,
		OutputSuffix:       "-flipped",
		ProgressBarMessage: "Flipping images",
		ProcessorFunc:      processFlipImage,
		ResultVerb:         "flipped",
	})
}

func processFlipImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewBimgImage(img)
		flippedImg, err := newImg.Flip()
		if err != nil {
			return nil, err
		}
		return flippedImg, nil
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
	Run: flipRun,
}

func init() {
	rootCmd.AddCommand(flipCmd)

	flipCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the images to flip")
	flipCmd.MarkFlagRequired("directory")
}
