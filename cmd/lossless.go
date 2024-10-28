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

func losslessRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fs := filesystem.NewFileSystem()
	files, err := fs.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	outputDir, err := fs.CreateSiblingDir(directory, "-lossless")
	if err != nil {
		log.Fatal(err)
	}

	var initialSize, finalSize uint64
	var skippedImages, losslessImages uint32

	resultBuilder := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(os.Stdout, len(files), concurrency, "Lossless images processing")

	params := processing.ProcessFilesParams{
		Files:       files,
		FS:          fs,
		OutputDir:   outputDir,
		ProgressBar: progressBar,
		ProcessorFunc: func(p processing.FileProcessingParams) error {
			stats := utils.NewImageProcessingStats(&initialSize, &finalSize, &skippedImages, &losslessImages)
			return processLosslessImage(ctx, p, stats)
		},
		Concurrency: concurrency,
	}
	processing.ProcessFiles(params)

	progressBar.Finish()

	totalImages := uint32(len(files))
	resultBuilder.SetTotalImages(totalImages).
		SetSkippedImages(skippedImages).
		SetProcessedImages(losslessImages).
		SetOutputDirectory(outputDir).
		SetInitialSize(float64(initialSize)).
		SetFinalSize(float64(finalSize))
	result := resultBuilder.Build()
	fmt.Println(result.PrintResults("lossless"))
}

func processLosslessImage(ctx context.Context, params processing.FileProcessingParams, stats *utils.ImageProcessingStats) error {
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
	compressedImg, err := newImg.LosslessCompress()
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	outputPath := utils.BuildOutputPath(params.OutputDir, params.File.Path)
	err = params.FS.WriteFile(outputPath, compressedImg)
	if err != nil {
		atomic.AddUint32(stats.SkippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.FinalSize, uint64(len(compressedImg)))
	atomic.AddUint32(stats.ProcessedImages, 1)
	return nil
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
