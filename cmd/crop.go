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
	"github.com/felipesimis/compactify-cli/pkg/validation"
	"github.com/h2non/bimg"
	"github.com/spf13/cobra"
)

var (
	gravity int
)

type CropParams struct {
	Width   int
	Height  int
	Gravity bimg.Gravity
}

type CropStats struct {
	initialSize   *uint64
	finalSize     *uint64
	skippedImages *uint32
	croppedImages *uint32
}

func cropRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dimensionValidation := &validation.DimensionsValidation{Width: width, Height: height}
	gravityValidation := &validation.GravityValidation{Gravity: bimg.Gravity(gravity)}
	validationComposite := validation.ValidationComposite{Validations: []validation.Validation{dimensionValidation, gravityValidation}}
	err := validationComposite.Validate()
	if err != nil {
		log.Fatal(err)
	}

	fs := filesystem.NewFileSystem()
	files, err := fs.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	outputDir, err := fs.CreateSiblingDir(directory, fmt.Sprintf("-cropped_%dx%d", width, height))
	if err != nil {
		log.Fatal(err)
	}

	var initialSize, finalSize uint64
	var skippedImages, croppedImages uint32

	resultBuilder := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(os.Stdout, len(files), concurrency, "Cropping images")

	params := processing.ProcessFilesParams{
		Files:       files,
		FS:          fs,
		OutputDir:   outputDir,
		ProgressBar: progressBar,
		ExtraParams: CropParams{Width: width, Height: height, Gravity: bimg.Gravity(gravity)},
		ProcessorFunc: func(p processing.FileProcessingParams) error {
			extraParams := p.ExtraParams.(CropParams)
			stats := &CropStats{
				initialSize:   &initialSize,
				finalSize:     &finalSize,
				skippedImages: &skippedImages,
				croppedImages: &croppedImages,
			}
			return cropImages(ctx, p, extraParams, stats)
		},
		Concurrency: concurrency,
	}
	processing.ProcessFiles(params)

	progressBar.Finish()

	totalImages := uint32(len(files))
	resultBuilder.SetTotalImages(totalImages).
		SetSkippedImages(skippedImages).
		SetProcessedImages(croppedImages).
		SetOutputDirectory(outputDir).
		SetInitialSize(float64(initialSize)).
		SetFinalSize(float64(finalSize))
	result := resultBuilder.Build()
	fmt.Println(result.PrintResults("cropped"))
}

func cropImages(ctx context.Context, params processing.FileProcessingParams, extraParams CropParams, stats *CropStats) error {
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
	croppedImg, err := newImg.Crop(extraParams.Width, extraParams.Height, extraParams.Gravity)
	if err != nil {
		atomic.AddUint32(stats.skippedImages, 1)
		return err
	}

	outputPath := utils.BuildOutputPath(params.OutputDir, params.File.Path)
	err = params.FS.WriteFile(outputPath, croppedImg)
	if err != nil {
		atomic.AddUint32(stats.skippedImages, 1)
		return err
	}

	atomic.AddUint64(stats.finalSize, uint64(len(croppedImg)))
	atomic.AddUint32(stats.croppedImages, 1)
	params.ProgressBar.Increment()
	return nil
}

var cropCmd = &cobra.Command{
	Use:     "crop",
	Aliases: []string{"cut"},
	Args:    cobra.NoArgs,
	Short:   "Crop an image to specified dimensions",
	Long: `Crop an image to a specific width and height.
This command allows you to change the dimensions of an image by cropping it, which can be useful for optimizing images for 
different uses, such as web, mobile, or print. You can specify the desired width and height, 
and the image will be cropped accordingly.`,
	Run: cropRun,
}

func init() {
	rootCmd.AddCommand(cropCmd)

	cropCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the image to crop")
	cropCmd.Flags().IntVarP(&width, "width", "w", 0, "Desired width of the image")
	cropCmd.Flags().IntVarP(&height, "height", "H", 0, "Desired height of the image")
	cropCmd.Flags().IntVarP(&gravity, "gravity", "g", int(bimg.GravityCentre), `Gravity to use when cropping the image. 
Available options:
  0 - Centre (default)
  1 - North
  2 - East
  3 - South
  4 - West
  5 - Smart`)

	cropCmd.MarkFlagRequired("directory")
	cropCmd.MarkFlagRequired("width")
	cropCmd.MarkFlagRequired("height")
}
