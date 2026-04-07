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
	const (
		Reset  = "\033[0m"
		Red    = "\033[1;31m"
		Green  = "\033[1;32m"
		Yellow = "\033[1;33m"
		Cyan   = "\033[1;36m"
		Bold   = "\033[1m"
	)

	var result strings.Builder
	fmt.Fprintf(&result, "\n%s============================================================%s\n", Bold, Reset)
	fmt.Fprintf(&result, "⏱️  Elapsed time: %s\n", r.elapsedTime)
	fmt.Fprintf(&result, "📁 Output directory: %s\n", r.outputDirectory)
	fmt.Fprintf(&result, "🖼️  Total images: %d\n", r.totalImages)

	if r.skippedImages > 0 {
		fmt.Fprintf(&result, "%s⏭️  Skipped images: %d%s\n", Yellow, r.skippedImages, Reset)
	} else {
		fmt.Fprintf(&result, "⏭️  Skipped images: %d\n", r.skippedImages)
	}

	processedLabel := strings.ToUpper(string(key[0])) + key[1:]
	fmt.Fprintf(&result, "%s✅ %s: %d%s\n", Green, processedLabel, r.processedImages, Reset)
	fmt.Fprintf(&result, "\n📦 Initial size: %.2f MB\n", r.initialSize)
	fmt.Fprintf(&result, "📦 Final size: %.2f MB\n", r.finalSize)
	fmt.Fprintf(&result, "%s💾 Size difference: %.2f MB%s\n", Cyan, r.sizeDifference, Reset)
	fmt.Fprintf(&result, "%s📉 Size difference percentage: %.2f%%%s\n", Cyan, r.sizeDifferencePercentage, Reset)

	if len(r.errors) > 0 {
		fmt.Fprintf(&result, "\n%s⚠️  Errors found during processing:%s\n", Red, Reset)
		for _, err := range r.errors {
			fmt.Fprintf(&result, "%s  ❌ %v%s\n", Red, err, Reset)
		}
	}
	fmt.Fprintf(&result, "%s============================================================%s", Bold, Reset)

	return result.String()
}
