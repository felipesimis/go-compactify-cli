package cmd

import (
	"context"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/image"
	"github.com/felipesimis/compactify-cli/internal/processing"
	"github.com/felipesimis/compactify-cli/internal/utils"
	"github.com/spf13/cobra"
)

func grayscaleRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fs := filesystem.NewFileSystem()
	RunOperation(OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
		InputDir:           directory,
		OutputSuffix:       "-grayscale",
		ProgressBarMessage: "Creating grayscale images",
		ProcessorFunc:      processGrayscaleImage,
		ResultVerb:         "grayscale images created",
	})
}

func processGrayscaleImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewBimgImage(img)
		return newImg.Grayscale()
	})
}

var grayscaleCmd = &cobra.Command{
	Use:     "grayscale",
	Aliases: []string{"gray", "bw"},
	Args:    cobra.NoArgs,
	Short:   "Convert images to grayscale",
	Long: `Convert images to grayscale.
This command allows you to convert an image to grayscale, removing all color information and leaving only shades of gray.
It can be useful for various image processing tasks, such as creating artistic effects or preparing images for printing.`,
	Run: grayscaleRun,
}

func init() {
	rootCmd.AddCommand(grayscaleCmd)

	grayscaleCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the images to grayscale")
	grayscaleCmd.MarkFlagRequired("directory")
}
