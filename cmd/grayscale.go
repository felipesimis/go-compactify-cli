package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/image"
	"github.com/felipesimis/compactify-cli/internal/processing"
	"github.com/felipesimis/compactify-cli/internal/utils"
	"github.com/felipesimis/compactify-cli/pkg/progress"
	"github.com/spf13/cobra"
)

type GrayscaleStats struct {
	initialSize     *uint64
	finalSize       *uint64
	skippedImages   *uint32
	grayscaleImages *uint32
}

func grayscaleRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fs := filesystem.NewFileSystem()
	files, err := fs.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	outputDir, err := fs.CreateSiblingDir(directory, "-grayscale")
	if err != nil {
		log.Fatal(err)
	}

	var initialSize, finalSize uint64
	var skippedImages, grayscaleImages uint32

	resultBuilder := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(os.Stdout, len(files), concurrency, "Grayscaling images")

	params := processing.ProcessFilesParams{
		Files:       files,
		FS:          fs,
		OutputDir:   outputDir,
		ProgressBar: progressBar,
		ProcessorFunc: func(p processing.FileProcessingParams) error {
			stats := &GrayscaleStats{
				initialSize:     &initialSize,
				finalSize:       &finalSize,
				skippedImages:   &skippedImages,
				grayscaleImages: &grayscaleImages,
			}
			return processGrayscaleImage(ctx, p, stats)
		},
		Concurrency: concurrency,
	}
	processing.ProcessFiles(params)

	progressBar.Finish()

	totalImages := uint32(len(files))
	resultBuilder.SetTotalImages(totalImages).
		SetSkippedImages(skippedImages).
		SetProcessedImages(grayscaleImages).
		SetOutputDirectory(outputDir).
		SetInitialSize(float64(initialSize)).
		SetFinalSize(float64(finalSize))
	result := resultBuilder.Build()
	fmt.Println(result.PrintResults("grayscale"))
}

func processGrayscaleImage(ctx context.Context, params processing.FileProcessingParams, stats *GrayscaleStats) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	img, err := params.FS.ReadFile(params.File.Path)
	if err != nil {
		atomic.AddUint32(stats.skippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.initialSize, uint64(params.File.Size))
	newImg := image.NewBimgImage(img)
	grayscaleImages, err := newImg.Grayscale()
	if err != nil {
		atomic.AddUint32(stats.skippedImages, 1)
		return err
	}

	outputPath := utils.BuildOutputPath(params.OutputDir, params.File.Path)
	err = params.FS.WriteFile(outputPath, grayscaleImages)
	if err != nil {
		atomic.AddUint32(stats.skippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.finalSize, uint64(len(grayscaleImages)))
	atomic.AddUint32(stats.grayscaleImages, 1)
	return nil
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
