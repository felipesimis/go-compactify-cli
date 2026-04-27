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

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreTotalImages() {
	suite.rb.SetTotalImages(10)
	suite.Equal(10, int(suite.rb.result.totalImages))
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreSkippedImages() {
	suite.rb.SetSkippedImages(5)
	suite.Equal(5, int(suite.rb.result.skippedImages))
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreProcessedImages() {
	suite.rb.SetProcessedImages(5)
	suite.Equal(5, int(suite.rb.result.processedImages))
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreOriginalBytes() {
	suite.rb.SetOriginalBytes(100)
	suite.Equal(int64(100), suite.rb.result.originalBytes)
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreProcessedBytes() {
	suite.rb.SetProcessedBytes(50)
	suite.Equal(int64(50), suite.rb.result.processedBytes)
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreOutputDirectory() {
	suite.rb.SetOutputDirectory("output")
	suite.Equal("output", suite.rb.result.outputDirectory)
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreErrors() {
	errors := []error{assert.AnError, assert.AnError}
	suite.rb.SetErrors(errors)
	suite.Equal(errors, suite.rb.result.errors)
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldBuildCorrectResult() {
	suite.rb.
		SetOriginalBytes(10485760). // 10 MB
		SetProcessedBytes(5242880). // 5 MB
		SetTotalImages(10).
		SetSkippedImages(3).
		SetProcessedImages(7).
		SetOutputDirectory("output").
		SetErrors([]error{assert.AnError})

	result := suite.rb.Build()

	suite.Equal(time.Second, result.elapsedTime)
	suite.Equal(int64(10485760), result.originalBytes)
	suite.Equal(int64(5242880), result.processedBytes)
	suite.Equal(int64(5242880), result.savedBytes)
	suite.Equal(uint32(10), result.totalImages)
	suite.Equal(uint32(3), result.skippedImages)
	suite.Equal(uint32(7), result.processedImages)
	suite.Equal(50.0, result.reductionRatio)
	suite.Equal("output", result.outputDirectory)
	suite.Equal([]error{assert.AnError}, result.errors)
	suite.mockTimeProvider.AssertExpectations(suite.T())
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldPrintFormattedResults() {
	tests := []struct {
		name            string
		skippedImages   uint32
		processedImages uint32
		errors          []error
		expected        []string
		notExpected     []string
	}{
		{
			name:            "without errors",
			skippedImages:   0,
			processedImages: 10,
			errors:          nil,
			expected: []string{
				"OPERATION", "IMPACT", "OUTPUT DIRECTORY",
				"10 images", "0", "10",
				"10.00 MB", "5.00 MB", "50.00%",
				"output",
			},
			notExpected: []string{
				"ERRORS DETECTED",
			},
		},
		{
			name:            "with errors",
			skippedImages:   3,
			processedImages: 7,
			errors: []error{
				fmt.Errorf("file 'fake.jpg': read error"),
				fmt.Errorf("permission denied"),
			},
			expected: []string{
				"2 ERRORS DETECTED",
				"fake.jpg",
				"read error",
				"permission denied",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.rb.
				SetOriginalBytes(10485760). // 10 MB
				SetProcessedBytes(5242880). // 5 MB
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
			for _, notExpectedText := range tt.notExpected {
				suite.NotContains(printedResult, notExpectedText)
			}
			suite.mockTimeProvider.AssertExpectations(suite.T())
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

func TestResultBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(ResultBuilderTestSuite))
}
