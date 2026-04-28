package cmd

import (
	"context"

	"github.com/felipesimis/go-compactify-cli/internal/filesystem"
	"github.com/felipesimis/go-compactify-cli/internal/image"
	"github.com/felipesimis/go-compactify-cli/internal/processing"
	"github.com/felipesimis/go-compactify-cli/internal/utils"
	"github.com/spf13/cobra"
)

func paletteRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	fs := filesystem.NewFileSystem()
	globalConfig := loadGlobalConfig(cmd)
	return RunOperation(globalConfig, OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
		OutputSuffix:       "-palette",
		ProgressBarMessage: "Enabling palette on images",
		ProcessorFunc:      processPaletteImage,
		ResultVerb:         "palette enabled",
	})
}

func processPaletteImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	return HandleImageProcessing(ctx, params, stats, func(proc image.ImageProcessor) ([]byte, error) {
		return proc.EnablePalette()
	})
}

var paletteCmd = &cobra.Command{
	Use:   "palette",
	Args:  cobra.NoArgs,
	Short: "Enable palette on images",
	Long: `Apply a color palette to images.
This command enables a color palette on the specified images, which can help reduce the file size by limiting the number of colors used. 
It is useful for optimizing images for web use, creating artistic effects, and ensuring compatibility with formats that require or benefit from a limited color palette.`,
	RunE: paletteRun,
}

func init() {
	rootCmd.AddCommand(paletteCmd)
}
