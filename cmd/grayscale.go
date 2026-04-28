package cmd

import (
	"context"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/spf13/cobra"
)

func grayscaleRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	fs := filesystem.NewFileSystem()
	globalConfig := loadGlobalConfig(cmd)
	return RunOperation(globalConfig, OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
		OutputSuffix:       "-grayscale",
		ProgressBarMessage: "Creating grayscale images",
		ProcessorFunc:      processGrayscaleImage,
		ResultVerb:         "grayscale images created",
	})
}

func processGrayscaleImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewProcessor(img)
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
	RunE: grayscaleRun,
}

func init() {
	rootCmd.AddCommand(grayscaleCmd)
}
