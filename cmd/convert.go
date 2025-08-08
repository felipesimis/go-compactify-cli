package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/image"
	"github.com/felipesimis/compactify-cli/internal/processing"
	"github.com/felipesimis/compactify-cli/internal/utils"
	"github.com/felipesimis/compactify-cli/pkg/validation"
	"github.com/spf13/cobra"
)

var format string

type ConvertParams struct {
	Format string
}

func convertRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dimensionValidation := &validation.FormatValidation{Format: format}
	err := dimensionValidation.Validate()
	if err != nil {
		log.Fatal(err)
	}

	fs := filesystem.NewFileSystem()
	RunOperation(OperationConfig{
		Ctx:                ctx,
		FileSystem:         fs,
		InputDir:           directory,
		OutputSuffix:       fmt.Sprintf("-converted.%s", format),
		ProgressBarMessage: "Converting images",
		ExtraParams:        ConvertParams{Format: format},
		ProcessorFunc:      processConvertImage,
		ResultVerb:         "converted",
	})
}

func processConvertImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
	extraParams := params.ExtraParams.(ConvertParams)
	return HandleImageProcessing(ctx, params, stats, func(img []byte) ([]byte, error) {
		newImg := image.NewBimgImage(img)
		return newImg.Convert(extraParams.Format)
	})
}

var convertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"conv"},
	Args:    cobra.NoArgs,
	Short:   "Convert images to a specified format",
	Long: `Convert images in a directory to a specified format.
This command allows you to change the format of images, which can be useful for optimizing images for 
different uses, such as web, mobile, or print. You can specify the desired format, 
and the images will be converted accordingly.`,
	Run: convertRun,
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the images to convert")
	convertCmd.Flags().StringVarP(&format, "format", "f", "", `Desired format of the images. Available options: webp, jpeg, png`)

	convertCmd.MarkFlagRequired("directory")
	convertCmd.MarkFlagRequired("format")
}
