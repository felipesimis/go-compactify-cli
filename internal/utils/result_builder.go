package utils

import (
	"fmt"
	"time"

	"github.com/felipesimis/compactify-cli/internal/ui"
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
	left := ui.Panel{
		Title: "OPERATION",
		Items: []ui.Item{
			{Label: "Elapsed time", Value: r.elapsedTime.Round(time.Millisecond).String(), IsHighlighted: false},
			{Label: "Total", Value: fmt.Sprintf("%d images", r.totalImages), IsHighlighted: false},
			{Label: "Skipped", Value: fmt.Sprintf("%d", r.skippedImages), IsHighlighted: false},
			{Label: "Processed", Value: fmt.Sprintf("%d", r.processedImages), IsHighlighted: false},
		},
	}

	toMB := func(b int64) float64 { return float64(b) / 1024 / 1024 }
	right := ui.Panel{
		Title: "IMPACT",
		Items: []ui.Item{
			{Label: "Original", Value: fmt.Sprintf("%.2f MB", toMB(r.originalBytes)), IsHighlighted: false},
			{Label: "After", Value: fmt.Sprintf("%.2f MB", toMB(r.processedBytes)), IsHighlighted: false},
			{Label: "", Value: ""},
			{Label: "Saved", Value: fmt.Sprintf("%.2f MB", toMB(r.savedBytes)), IsHighlighted: true},
			{Label: "Reduction", Value: fmt.Sprintf("%.2f%%", r.reductionRatio), IsHighlighted: true},
		},
	}

	dashboard := ui.RenderDashboard(left, right, "OUTPUT DIRECTORY", fmt.Sprintf("📂 %s", r.outputDirectory))
	errors := ui.RenderErrorList(r.errors)

	return "\n" + dashboard + errors + "\n"
}
