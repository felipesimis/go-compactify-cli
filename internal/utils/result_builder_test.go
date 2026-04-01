package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TimeMock struct {
	mock.Mock
}

func (tm *TimeMock) Now() time.Time {
	args := tm.Called()
	return args.Get(0).(time.Time)
}

func (tm *TimeMock) Since(t time.Time) time.Duration {
	args := tm.Called(t)
	return args.Get(0).(time.Duration)
}

func mockTimeProvider() *TimeMock {
	timeMock := new(TimeMock)
	startTime := time.Now()
	timeMock.On("Now").Return(startTime)
	timeMock.On("Since", startTime).Return(time.Second)
	return timeMock
}

func TestResultBuilder_SetTotalImages(t *testing.T) {
	rb := NewResultBuilder(mockTimeProvider())
	rb.SetTotalImages(10)
	assert.Equal(t, 10, int(rb.result.totalImages))
}

func TestResultBuilder_SetSkippedImages(t *testing.T) {
	rb := NewResultBuilder(mockTimeProvider())
	rb.SetSkippedImages(5)
	assert.Equal(t, 5, int(rb.result.skippedImages))
}

func TestResultBuilder_SetProcessedImages(t *testing.T) {
	rb := NewResultBuilder(mockTimeProvider())
	rb.SetProcessedImages(5)
	assert.Equal(t, 5, int(rb.result.processedImages))
}

func TestResultBuilder_SetInitialSize(t *testing.T) {
	rb := NewResultBuilder(mockTimeProvider())
	rb.SetInitialSize(100)
	assert.Equal(t, 100.0, rb.result.initialSize)
}

func TestResultBuilder_SetFinalSize(t *testing.T) {
	rb := NewResultBuilder(mockTimeProvider())
	rb.SetFinalSize(50)
	assert.Equal(t, 50.0, rb.result.finalSize)
}

func TestResultBuilder_SetOutputDirectory(t *testing.T) {
	rb := NewResultBuilder(mockTimeProvider())
	rb.SetOutputDirectory("output")
	assert.Equal(t, "output", rb.result.outputDirectory)
}

func TestResultBuilder_SetErrors(t *testing.T) {
	rb := NewResultBuilder(mockTimeProvider())
	errors := []error{assert.AnError, assert.AnError}
	rb.SetErrors(errors)
	assert.Equal(t, errors, rb.result.errors)
}

func TestResultBuilder_Build(t *testing.T) {
	timeMock := mockTimeProvider()
	rb := NewResultBuilder(timeMock).
		SetInitialSize(10485760). // 10 MB
		SetFinalSize(5242880).    // 5 MB
		SetTotalImages(10).
		SetSkippedImages(3).
		SetProcessedImages(7).
		SetOutputDirectory("output").
		SetErrors([]error{assert.AnError})
	result := rb.Build()

	assert.Equal(t, time.Second, result.elapsedTime)
	assert.Equal(t, 10.0, result.initialSize)
	assert.Equal(t, 5.0, result.finalSize)
	assert.Equal(t, -5.0, result.sizeDifference)
	assert.Equal(t, 10, int(result.totalImages))
	assert.Equal(t, 3, int(result.skippedImages))
	assert.Equal(t, 7, int(result.processedImages))
	assert.Equal(t, 50.0, result.sizeDifferencePercentage)
	assert.Equal(t, "output", result.outputDirectory)
	assert.Equal(t, []error{assert.AnError}, result.errors)
	timeMock.AssertExpectations(t)
}

func TestResultBuilder_Result_PrintResults(t *testing.T) {
	tests := []struct {
		name     string
		errors   []error
		expected string
	}{
		{
			name:   "Without errors",
			errors: nil,
			expected: `Elapsed time: 1s
Total images: 10
Skipped images: 3
Resized images: 7
Initial size: 10.00 MB
Final size: 5.00 MB
Size difference: -5.00 MB
Size difference percentage: 50.00%
Output directory: output`,
		},
		{
			name:   "With errors",
			errors: []error{fmt.Errorf("file 'fake.jpg': read error")},
			expected: `Elapsed time: 1s
Total images: 10
Skipped images: 3
Resized images: 7
Initial size: 10.00 MB
Final size: 5.00 MB
Size difference: -5.00 MB
Size difference percentage: 50.00%
Output directory: output
Errors found during processing:
  ❌ file 'fake.jpg': read error
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeMock := mockTimeProvider()

			rb := NewResultBuilder(timeMock).
				SetInitialSize(10485760). // 10 MB
				SetFinalSize(5242880).    // 5 MB
				SetTotalImages(10).
				SetSkippedImages(3).
				SetProcessedImages(7).
				SetOutputDirectory("output")

			if tt.errors != nil {
				rb.SetErrors(tt.errors)
			}
			result := rb.Build()
			printedResult := result.PrintResults("resized")

			assert.Equal(t, tt.expected, printedResult)
			timeMock.AssertExpectations(t)
		})
	}
}

func TestRealTimeProvider_Now(t *testing.T) {
	rtp := RealTimeProvider{}
	now := time.Now()
	rtpNow := rtp.Now()

	assert.WithinDuration(t, now, rtpNow, time.Second)
}

func TestRealTimeProvider_Since(t *testing.T) {
	rtp := RealTimeProvider{}
	startTime := time.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := rtp.Since(startTime)

	assert.InDelta(t, 100*time.Millisecond, elapsed, float64(time.Millisecond))
}
