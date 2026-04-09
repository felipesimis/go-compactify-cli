package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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

func newMockTimeProvider() *TimeMock {
	timeMock := new(TimeMock)
	startTime := time.Now()
	timeMock.On("Now").Return(startTime)
	timeMock.On("Since", startTime).Return(time.Second)
	return timeMock
}

type ResultBuilderTestSuite struct {
	suite.Suite
	mockTimeProvider *TimeMock
	rb               *ResultBuilder
}

func (suite *ResultBuilderTestSuite) SetupTest() {
	suite.mockTimeProvider = newMockTimeProvider()
	suite.rb = NewResultBuilder(suite.mockTimeProvider)
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_SetTotalImages() {
	suite.rb.SetTotalImages(10)
	suite.Equal(10, int(suite.rb.result.totalImages))
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_SetSkippedImages() {
	suite.rb.SetSkippedImages(5)
	suite.Equal(5, int(suite.rb.result.skippedImages))
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_SetProcessedImages() {
	suite.rb.SetProcessedImages(5)
	suite.Equal(5, int(suite.rb.result.processedImages))
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_SetInitialSize() {
	suite.rb.SetInitialSize(100)
	suite.Equal(100.0, suite.rb.result.initialSize)
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_SetFinalSize() {
	suite.rb.SetFinalSize(50)
	suite.Equal(50.0, suite.rb.result.finalSize)
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_SetOutputDirectory() {
	suite.rb.SetOutputDirectory("output")
	suite.Equal("output", suite.rb.result.outputDirectory)
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_SetErrors() {
	errors := []error{assert.AnError, assert.AnError}
	suite.rb.SetErrors(errors)
	suite.Equal(errors, suite.rb.result.errors)
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_Build() {
	suite.rb.
		SetInitialSize(10485760). // 10 MB
		SetFinalSize(5242880).    // 5 MB
		SetTotalImages(10).
		SetSkippedImages(3).
		SetProcessedImages(7).
		SetOutputDirectory("output").
		SetErrors([]error{assert.AnError})
	result := suite.rb.Build()

	suite.Equal(time.Second, result.elapsedTime)
	suite.Equal(10.0, result.initialSize)
	suite.Equal(5.0, result.finalSize)
	suite.Equal(-5.0, result.sizeDifference)
	suite.Equal(10, int(result.totalImages))
	suite.Equal(3, int(result.skippedImages))
	suite.Equal(7, int(result.processedImages))
	suite.Equal(50.0, result.sizeDifferencePercentage)
	suite.Equal("output", result.outputDirectory)
	suite.Equal([]error{assert.AnError}, result.errors)
	suite.mockTimeProvider.AssertExpectations(suite.T())
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_Result_PrintResults() {
	tests := []struct {
		name            string
		skippedImages   uint32
		processedImages uint32
		errors          []error
		expected        []string
	}{
		{
			name:            "Without errors",
			skippedImages:   3,
			processedImages: 7,
			errors:          nil,
			expected: []string{
				"Elapsed time: 1s",
				"Total images: 10",
				"Skipped images: 3",
				"Resized: 7",
				"Initial size: 10.00 MB",
				"Final size: 5.00 MB",
				"Size difference: -5.00 MB",
				"Size difference percentage: 50.00%",
				"Output directory: output",
			},
		},
		{
			name:            "With errors",
			skippedImages:   3,
			processedImages: 7,
			errors:          []error{fmt.Errorf("file 'fake.jpg': read error")},
			expected: []string{
				"Elapsed time: 1s",
				"Total images: 10",
				"Skipped images: 3",
				"Resized: 7",
				"Initial size: 10.00 MB",
				"Final size: 5.00 MB",
				"Size difference: -5.00 MB",
				"Size difference percentage: 50.00%",
				"Output directory: output",
				"Errors found during processing:",
				"  ❌ file 'fake.jpg': read error",
			},
		},
		{
			name:            "Without skipped images",
			skippedImages:   0,
			processedImages: 10,
			errors:          nil,
			expected: []string{
				"Elapsed time: 1s",
				"Total images: 10",
				"⏭️  Skipped images: 0",
				"Resized: 10",
				"Initial size: 10.00 MB",
				"Final size: 5.00 MB",
				"Size difference: -5.00 MB",
				"Size difference percentage: 50.00%",
				"Output directory: output",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.rb.
				SetInitialSize(10485760). // 10 MB
				SetFinalSize(5242880).    // 5 MB
				SetTotalImages(10).
				SetSkippedImages(tt.skippedImages).
				SetProcessedImages(tt.processedImages).
				SetOutputDirectory("output")

			if tt.errors != nil {
				suite.rb.SetErrors(tt.errors)
			}
			result := suite.rb.Build()
			printedResult := result.PrintResults("resized")

			for _, expectedText := range tt.expected {
				suite.Contains(printedResult, expectedText)
			}
			suite.mockTimeProvider.AssertExpectations(suite.T())
		})
	}
}

func (suite *ResultBuilderTestSuite) TestRealTimeProvider_Now() {
	rtp := RealTimeProvider{}
	now := time.Now()
	rtpNow := rtp.Now()

	suite.WithinDuration(now, rtpNow, time.Second)
}

func (suite *ResultBuilderTestSuite) TestRealTimeProvider_Since() {
	rtp := RealTimeProvider{}
	startTime := time.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := rtp.Since(startTime)

	suite.InDelta(100*time.Millisecond, elapsed, float64(time.Millisecond))
}

func TestResultBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(ResultBuilderTestSuite))
}
