package cmd

import (
	"context"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/image"
	"github.com/felipesimis/compactify-cli/internal/processing"
	"github.com/felipesimis/compactify-cli/internal/utils"
	"github.com/spf13/cobra"
)

func losslessRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fs := filesystem.NewFileSystem()
	RunOperation(OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
		InputDir:           directory,
		OutputSuffix:       "-lossless",
		ProgressBarMessage: "Applying lossless compression",
		ProcessorFunc:      processLosslessImage,
		ResultVerb:         "lossless compressed",
	})
}

func processLosslessImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewBimgImage(img)
		compressedImg, err := newImg.LosslessCompress()
		if err != nil {
			return nil, err
		}
		return compressedImg, nil
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
	Run: losslessRun,
}

func init() {
	rootCmd.AddCommand(losslessCmd)

	losslessCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the images to compress")
	losslessCmd.MarkFlagRequired("directory")
}
