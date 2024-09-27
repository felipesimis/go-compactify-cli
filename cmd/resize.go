/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"

	"github.com/felipesimis/compactify-cli/internal/filesystem"
	"github.com/felipesimis/compactify-cli/internal/image"
	"github.com/felipesimis/compactify-cli/internal/utils"
	"github.com/felipesimis/compactify-cli/pkg/progress"
	"github.com/felipesimis/compactify-cli/pkg/validation"
	"github.com/spf13/cobra"
)

var (
	directory string
	width     int
	height    int
)

type ProcessFileParams struct {
	fileInfo      filesystem.FileInfo
	fs            filesystem.FileSystem
	outputDir     string
	initialSize   *uint64
	finalSize     *uint64
	skippedImages *uint32
	resizedImages *uint32
	errChan       chan error
	width         int
	height        int
	progressBar   *progress.ProgressBar
}

func resizeRun(cmd *cobra.Command, args []string) {
	fs := filesystem.NewFileSystem()
	files, err := fs.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	outputDir, err := fs.CreateSiblingDir(directory, "-resized")
	if err != nil {
		log.Fatal(err)
	}

	sem := make(chan struct{}, concurrency)
	errChan := make(chan error, len(files))
	var wg sync.WaitGroup
	var initialSize, finalSize uint64
	var skippedImages, resizedImages uint32

	rb := utils.NewResultBuilder(utils.RealTimeProvider{})
	progressBar := progress.NewProgressBar(os.Stdout, len(files), "Resizing images")

	for _, fileInfo := range files {
		wg.Add(1)
		sem <- struct{}{}
		go func(fileInfo filesystem.FileInfo) {
			defer func() {
				<-sem
				wg.Done()
			}()
			params := ProcessFileParams{
				fileInfo:      fileInfo,
				fs:            fs,
				outputDir:     outputDir,
				initialSize:   &initialSize,
				finalSize:     &finalSize,
				skippedImages: &skippedImages,
				resizedImages: &resizedImages,
				errChan:       errChan,
				width:         width,
				height:        height,
				progressBar:   progressBar,
			}
			processFile(params)
		}(fileInfo)
	}

	wg.Wait()
	close(errChan)

	progressBar.Finish()

	for err := range errChan {
		log.Println(err)
	}

	totalImages := uint32(len(files))
	rb.SetTotalImages(totalImages).
		SetSkippedImages(skippedImages).
		SetProcessedImages(resizedImages).
		SetOutputDirectory(outputDir).
		SetInitialSize(float64(initialSize)).
		SetFinalSize(float64(finalSize))
	result := rb.Build()
	fmt.Println(result.PrintResults("resized"))
}

func processFile(params ProcessFileParams) {
	dimensionValidation := &validation.DimensionsValidation{Width: params.width, Height: params.height}
	err := dimensionValidation.Validate()
	if err != nil {
		log.Fatal(err)
	}

	img, err := params.fs.ReadFile(params.fileInfo.Path)
	if err != nil {
		params.errChan <- err
		atomic.AddUint32(params.skippedImages, 1)
		params.progressBar.Increment()
		return
	}

	atomic.AddUint64(params.initialSize, uint64(params.fileInfo.Size))
	newImg := image.NewBimgImage(img)
	resizedImg, err := newImg.Resize(params.width, params.height)
	if err != nil {
		params.errChan <- err
		atomic.AddUint32(params.skippedImages, 1)
		params.progressBar.Increment()
		return
	}

	outputPath := utils.BuildOutputPath(params.outputDir, params.fileInfo.Path)
	err = params.fs.WriteFile(outputPath, resizedImg)
	if err != nil {
		params.errChan <- err
		atomic.AddUint32(params.skippedImages, 1)
		params.progressBar.Increment()
		return
	}

	atomic.AddUint64(params.finalSize, uint64(len(resizedImg)))
	atomic.AddUint32(params.resizedImages, 1)
	params.progressBar.Increment()
}

var resizeCmd = &cobra.Command{
	Use:   "resize",
	Args:  cobra.NoArgs,
	Short: "Resize an image to specified dimensions",
	Long: `Resize an image to a specific width and height.
This command allows you to change the dimensions of an image, which can be useful for optimizing images for 
different uses, such as web, mobile, or print. You can specify the desired width and height, 
and the image will be resized accordingly.`,
	Run: resizeRun,
}

func init() {
	rootCmd.AddCommand(resizeCmd)

	resizeCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the image to resize")
	resizeCmd.Flags().IntVarP(&width, "width", "w", 0, "Desired width of the image")
	resizeCmd.Flags().IntVarP(&height, "height", "H", 0, "Desired height of the image")
	resizeCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 20, "Number of concurrent operations")

	resizeCmd.MarkFlagRequired("directory")
	resizeCmd.MarkFlagRequired("width")
	resizeCmd.MarkFlagRequired("height")
}
