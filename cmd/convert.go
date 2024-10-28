package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/image"
	"github.com/felipesimis/compactify-cli/internal/processing"
	"github.com/felipesimis/compactify-cli/internal/utils"
	"github.com/felipesimis/compactify-cli/pkg/progress"
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
	files, err := fs.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	outputDir, err := fs.CreateSiblingDir(directory, "-converted")
	if err != nil {
		log.Fatal(err)
	}

	var initialSize, finalSize uint64
	var skippedImages, convertedImages uint32

	resultBuilder := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(os.Stdout, len(files), concurrency, "Converting images")

	params := processing.ProcessFilesParams{
		Files:       files,
		FS:          fs,
		OutputDir:   outputDir,
		ProgressBar: progressBar,
		ExtraParams: ConvertParams{Format: format},
		ProcessorFunc: func(p processing.FileProcessingParams) error {
			extraParams := p.ExtraParams.(ConvertParams)
			stats := utils.NewImageProcessingStats(&initialSize, &finalSize, &skippedImages, &convertedImages)
			return processConvertImage(ctx, p, extraParams, stats)
		},
		Concurrency: concurrency,
	}
	processing.ProcessFiles(params)

	progressBar.Finish()

	totalImages := uint32(len(files))
	resultBuilder.SetTotalImages(totalImages).
		SetSkippedImages(skippedImages).
		SetProcessedImages(convertedImages).
		SetOutputDirectory(outputDir).
		SetInitialSize(float64(initialSize)).
		SetFinalSize(float64(finalSize))
	result := resultBuilder.Build()
	fmt.Println(result.PrintResults("converted"))
}

func processConvertImage(ctx context.Context, params processing.FileProcessingParams, extraParams ConvertParams, stats *utils.ImageProcessingStats) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	img, err := params.FS.ReadFile(params.File.Path)
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.InitialSize, uint64(params.File.Size))

	newImg := image.NewBimgImage(img)
	convertedImg, err := newImg.Convert(extraParams.Format)
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	outputPath := utils.BuildOutputPath(params.OutputDir, params.File.Path)
	newName := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + "." + extraParams.Format
	err = params.FS.WriteFile(newName, convertedImg)
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.FinalSize, uint64(len(convertedImg)))
	atomic.AddUint32(stats.ProcessedImages, 1)
	params.ProgressBar.Increment()
	return nil
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
