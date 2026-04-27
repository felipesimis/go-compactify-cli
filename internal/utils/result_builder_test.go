package utils

import (
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
	suite.Equal(10, int(suite.rb.result.TotalImages))
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreSkippedImages() {
	suite.rb.SetSkippedImages(5)
	suite.Equal(5, int(suite.rb.result.SkippedImages))
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreProcessedImages() {
	suite.rb.SetProcessedImages(5)
	suite.Equal(5, int(suite.rb.result.ProcessedImages))
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreOriginalBytes() {
	suite.rb.SetOriginalBytes(100)
	suite.Equal(int64(100), suite.rb.result.OriginalBytes)
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreProcessedBytes() {
	suite.rb.SetProcessedBytes(50)
	suite.Equal(int64(50), suite.rb.result.ProcessedBytes)
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreOutputDirectory() {
	suite.rb.SetOutputDirectory("output")
	suite.Equal("output", suite.rb.result.OutputDirectory)
}

func (suite *ResultBuilderTestSuite) TestResultBuilder_ShouldStoreErrors() {
	errors := []error{assert.AnError, assert.AnError}
	suite.rb.SetErrors(errors)
	suite.Equal(errors, suite.rb.result.Errors)
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

	suite.Equal(time.Second, result.ElapsedTime)
	suite.Equal(int64(10485760), result.OriginalBytes)
	suite.Equal(int64(5242880), result.ProcessedBytes)
	suite.Equal(int64(5242880), result.SavedBytes)
	suite.Equal(uint32(10), result.TotalImages)
	suite.Equal(uint32(3), result.SkippedImages)
	suite.Equal(uint32(7), result.ProcessedImages)
	suite.Equal(50.0, result.ReductionRatio)
	suite.Equal("output", result.OutputDirectory)
	suite.Equal([]error{assert.AnError}, result.Errors)
	suite.mockTimeProvider.AssertExpectations(suite.T())
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
