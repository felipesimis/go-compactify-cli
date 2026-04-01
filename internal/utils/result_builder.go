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
	startTime                time.Time
	elapsedTime              time.Duration
	totalImages              uint32
	skippedImages            uint32
	processedImages          uint32
	initialSize              float64
	finalSize                float64
	sizeDifference           float64
	sizeDifferencePercentage float64
	outputDirectory          string
	errors                   []error
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

func (rb *ResultBuilder) SetErrors(errs []error) *ResultBuilder {
	rb.result.errors = errs
	return rb
}

func (rb *ResultBuilder) Build() *Result {
	elapsedTime := rb.timeProvider.Since(rb.result.startTime)
	initialSizeMB := rb.result.initialSize / bytesInMb
	finalSizeMB := rb.result.finalSize / bytesInMb
	sizeDifferenceMB := initialSizeMB - finalSizeMB
	sizeDifferencePercentage := (sizeDifferenceMB / initialSizeMB) * 100

	if finalSizeMB < initialSizeMB {
		sizeDifferenceMB = -sizeDifferenceMB
	}

	return &Result{
		elapsedTime:              elapsedTime,
		totalImages:              rb.result.totalImages,
		skippedImages:            rb.result.skippedImages,
		processedImages:          rb.result.processedImages,
		initialSize:              initialSizeMB,
		finalSize:                finalSizeMB,
		sizeDifference:           sizeDifferenceMB,
		sizeDifferencePercentage: sizeDifferencePercentage,
		outputDirectory:          rb.result.outputDirectory,
		errors:                   rb.result.errors,
	}
}

func (r *Result) PrintResults(key string) string {
	var result strings.Builder
	fmt.Fprintf(&result, "Elapsed time: %s\n", r.elapsedTime)
	fmt.Fprintf(&result, "Total images: %d\n", r.totalImages)
	fmt.Fprintf(&result, "Skipped images: %d\n", r.skippedImages)
	fmt.Fprintf(&result, "%s images: %d\n", strings.ToUpper(string(key[0]))+key[1:], r.processedImages)
	fmt.Fprintf(&result, "Initial size: %.2f MB\n", r.initialSize)
	fmt.Fprintf(&result, "Final size: %.2f MB\n", r.finalSize)
	fmt.Fprintf(&result, "Size difference: %.2f MB\n", r.sizeDifference)
	fmt.Fprintf(&result, "Size difference percentage: %.2f%%\n", r.sizeDifferencePercentage)
	fmt.Fprintf(&result, "Output directory: %s", r.outputDirectory)

	if len(r.errors) > 0 {
		result.WriteString("\nErrors found during processing:\n")
		for _, err := range r.errors {
			fmt.Fprintf(&result, "  ❌ %v\n", err)
		}
	}
	return result.String()
}
