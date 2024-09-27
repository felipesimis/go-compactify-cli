package utils

import (
	"fmt"
	"strings"
	"time"
)

const (
	bytesInMb = 1024 * 1024
)

type TimeProvider interface {
	Now() time.Time
	Since(t time.Time) time.Duration
}

type RealTimeProvider struct{}

func (RealTimeProvider) Now() time.Time {
	return time.Now()
}

func (RealTimeProvider) Since(t time.Time) time.Duration {
	return time.Since(t)
}

type Result struct {
	startTime           time.Time
	elapsedTime         time.Duration
	totalImages         uint32
	skippedImages       uint32
	processedImages     uint32
	initialSize         float64
	finalSize           float64
	savedSize           float64
	savedSizePercentage float64
	outputDirectory     string
}

type ResultBuilder struct {
	result       *Result
	timeProvider TimeProvider
}

func NewResultBuilder(tp TimeProvider) *ResultBuilder {
	return &ResultBuilder{
		result: &Result{
			startTime: tp.Now(),
		},
		timeProvider: tp,
	}
}

func (rb *ResultBuilder) SetTotalImages(total uint32) *ResultBuilder {
	rb.result.totalImages = total
	return rb
}

func (rb *ResultBuilder) SetSkippedImages(skipped uint32) *ResultBuilder {
	rb.result.skippedImages = skipped
	return rb
}

func (rb *ResultBuilder) SetProcessedImages(resized uint32) *ResultBuilder {
	rb.result.processedImages = resized
	return rb
}

func (rb *ResultBuilder) SetInitialSize(size float64) *ResultBuilder {
	rb.result.initialSize = size
	return rb
}

func (rb *ResultBuilder) SetFinalSize(size float64) *ResultBuilder {
	rb.result.finalSize = size
	return rb
}

func (rb *ResultBuilder) SetOutputDirectory(directory string) *ResultBuilder {
	rb.result.outputDirectory = directory
	return rb
}

func (rb *ResultBuilder) Build() *Result {
	elapsedTime := rb.timeProvider.Since(rb.result.startTime)
	initialSizeMB := rb.result.initialSize / bytesInMb
	finalSizeMB := rb.result.finalSize / bytesInMb
	savedSizeMB := initialSizeMB - finalSizeMB
	savedSizePercentage := (savedSizeMB / initialSizeMB) * 100

	return &Result{
		elapsedTime:         elapsedTime,
		totalImages:         rb.result.totalImages,
		skippedImages:       rb.result.skippedImages,
		processedImages:     rb.result.processedImages,
		initialSize:         initialSizeMB,
		finalSize:           finalSizeMB,
		savedSize:           savedSizeMB,
		savedSizePercentage: savedSizePercentage,
		outputDirectory:     rb.result.outputDirectory,
	}
}

func (r *Result) PrintResults(key string) string {
	result := fmt.Sprintf("Elapsed time: %s\n", r.elapsedTime)
	result += fmt.Sprintf("Total images: %d\n", r.totalImages)
	result += fmt.Sprintf("Skipped images: %d\n", r.skippedImages)
	result += fmt.Sprintf("%s images: %d\n", strings.ToUpper(string(key[0]))+key[1:], r.processedImages)
	result += fmt.Sprintf("Initial size: %.2f MB\n", r.initialSize)
	result += fmt.Sprintf("Final size: %.2f MB\n", r.finalSize)
	result += fmt.Sprintf("Saved size: %.2f MB\n", r.savedSize)
	result += fmt.Sprintf("Saved size percentage: %.2f%%\n", r.savedSizePercentage)
	result += fmt.Sprintf("Output directory: %s", r.outputDirectory)
	return result
}
