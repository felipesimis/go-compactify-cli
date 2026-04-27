package utils

import (
	"time"
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
	StartTime       time.Time
	ElapsedTime     time.Duration
	TotalImages     uint32
	SkippedImages   uint32
	ProcessedImages uint32
	OriginalBytes   int64
	ProcessedBytes  int64
	SavedBytes      int64
	ReductionRatio  float64
	OutputDirectory string
	Errors          []error
}

type ResultBuilder struct {
	result       *Result
	timeProvider TimeProvider
}

func NewResultBuilder(tp TimeProvider) *ResultBuilder {
	return &ResultBuilder{
		result: &Result{
			StartTime: tp.Now(),
		},
		timeProvider: tp,
	}
}

func (rb *ResultBuilder) SetTotalImages(total uint32) *ResultBuilder {
	rb.result.TotalImages = total
	return rb
}

func (rb *ResultBuilder) SetSkippedImages(skipped uint32) *ResultBuilder {
	rb.result.SkippedImages = skipped
	return rb
}

func (rb *ResultBuilder) SetProcessedImages(resized uint32) *ResultBuilder {
	rb.result.ProcessedImages = resized
	return rb
}

func (rb *ResultBuilder) SetOriginalBytes(size uint64) *ResultBuilder {
	rb.result.OriginalBytes = int64(size)
	return rb
}

func (rb *ResultBuilder) SetProcessedBytes(size uint64) *ResultBuilder {
	rb.result.ProcessedBytes = int64(size)
	return rb
}

func (rb *ResultBuilder) SetOutputDirectory(directory string) *ResultBuilder {
	rb.result.OutputDirectory = directory
	return rb
}

func (rb *ResultBuilder) SetErrors(errs []error) *ResultBuilder {
	rb.result.Errors = errs
	return rb
}

func (rb *ResultBuilder) Build() *Result {
	elapsedTime := rb.timeProvider.Since(rb.result.StartTime)
	originalBytes := rb.result.OriginalBytes
	processedBytes := rb.result.ProcessedBytes
	savedBytes := originalBytes - processedBytes
	var reductionRatio float64

	if originalBytes > 0 {
		reductionRatio = (float64(savedBytes) / float64(originalBytes)) * 100
	}

	return &Result{
		ElapsedTime:     elapsedTime,
		TotalImages:     rb.result.TotalImages,
		SkippedImages:   rb.result.SkippedImages,
		ProcessedImages: rb.result.ProcessedImages,
		OriginalBytes:   originalBytes,
		ProcessedBytes:  processedBytes,
		SavedBytes:      savedBytes,
		ReductionRatio:  reductionRatio,
		OutputDirectory: rb.result.OutputDirectory,
		Errors:          rb.result.Errors,
	}
}
