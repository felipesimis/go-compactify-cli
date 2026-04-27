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
	startTime       time.Time
	elapsedTime     time.Duration
	totalImages     uint32
	skippedImages   uint32
	processedImages uint32
	originalBytes   int64
	processedBytes  int64
	savedBytes      int64
	reductionRatio  float64
	outputDirectory string
	errors          []error
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

func (rb *ResultBuilder) SetOriginalBytes(size uint64) *ResultBuilder {
	rb.result.originalBytes = int64(size)
	return rb
}

func (rb *ResultBuilder) SetProcessedBytes(size uint64) *ResultBuilder {
	rb.result.processedBytes = int64(size)
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
	originalBytes := rb.result.originalBytes
	processedBytes := rb.result.processedBytes
	savedBytes := originalBytes - processedBytes
	var reductionRatio float64

	if originalBytes > 0 {
		reductionRatio = (float64(savedBytes) / float64(originalBytes)) * 100
	}

	return &Result{
		elapsedTime:     elapsedTime,
		totalImages:     rb.result.totalImages,
		skippedImages:   rb.result.skippedImages,
		processedImages: rb.result.processedImages,
		originalBytes:   originalBytes,
		processedBytes:  processedBytes,
		savedBytes:      savedBytes,
		reductionRatio:  reductionRatio,
		outputDirectory: rb.result.outputDirectory,
		errors:          rb.result.errors,
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
	fmt.Fprintf(&result, "\n📦 Initial size: %.2f MB\n", float64(r.originalBytes)/bytesInMb)
	fmt.Fprintf(&result, "📦 Final size: %.2f MB\n", float64(r.processedBytes)/bytesInMb)
	fmt.Fprintf(&result, "%s💾 Size difference: %.2f MB%s\n", Cyan, float64(r.savedBytes)/bytesInMb, Reset)
	fmt.Fprintf(&result, "%s📉 Size difference percentage: %.2f%%%s\n", Cyan, r.reductionRatio, Reset)

	if len(r.errors) > 0 {
		fmt.Fprintf(&result, "\n%s⚠️  Errors found during processing:%s\n", Red, Reset)
		for _, err := range r.errors {
			fmt.Fprintf(&result, "%s  ❌ %v%s\n", Red, err, Reset)
		}
	}
	fmt.Fprintf(&result, "%s============================================================%s", Bold, Reset)

	return result.String()
}
